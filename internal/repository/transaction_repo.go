/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package repository

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
	"unit-service/internal/model/dao"
	"unit-service/internal/model/dto"
	"unit-service/internal/store"
	"unit-service/logger"
	"unit-service/metrics"

	"github.com/ClickHouse/clickhouse-go/v2"
)

const (
	batchSize         = 1000
	batchChanBuffSize = 3 * batchSize
	batchFlushTimeout = 300 * time.Millisecond
	batchPushTimeout  = 3 * time.Second
)

type TransactionRepo interface {
	RunBatchWriter(ctx context.Context) error
	PutBatch(ctx context.Context, transaction *dao.TrProc) error
	PutTransaction(transaction *dao.TranProc) error
	GetConnValid() bool
	SetConnValid(valid bool)
	ConnRecovery(ctx context.Context) error
}

type transactionRepo struct {
	conn        clickhouse.Conn
	store       store.TransactionStore
	batchCh     chan dao.TranProc
	isConnValid atomic.Bool
}

func NewTransactionRepo(store store.TransactionStore) (TransactionRepo, error) {
	conn := store.Conn()

	if conn == nil {
		return nil, errors.New("clickhouse connection is nil")
	}

	return &transactionRepo{
		conn:    conn,
		store:   store,
		batchCh: make(chan dao.TranProc, batchChanBuffSize),
	}, nil
}

func (repo *transactionRepo) PutTransaction(transaction *dao.Transaction) error {
	if repo.conn == nil {
		return errors.New("conn is nil")
	}

	query := getRepoInsQuery(dao.Transaction{})

	if err := repo.conn.Exec(context.Background(), query, gettrFields(transaction)...); err != nil {
		return err
	}

	return nil
}

func (repo *transactionRepo) PutBatch(ctx context.Context, transaction *dao.TranProc) error {
	if transaction == nil {
		return errors.New("transaction is nil")
	}

	select {
	case repo.batchCh <- *transaction:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (repo *transactionRepo) RunBatchWriter(ctx context.Context) error {
	ticker := time.NewTicker(batchFlushTimeout)
	defer ticker.Stop()

	batch := make([]dao.TranProc, 0, batchSize)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case transaction := <-repo.batchCh:
			batch = append(batch, transaction)
			if len(batch) >= batchSize {
				batch = repo.flushBatch(ctx, batch)
			}
		case <-ticker.C:
			batch = repo.flushBatch(ctx, batch)
		}
	}
}

func (repo *transactionRepo) runPush(ctx context.Context, buff []dao.TranProc) error {
	batch, err := repo.conn.PrepareBatch(ctx, dao.GetRepoInsQuery())
	if err != nil {
		return err
	}

	for _, tr := range buff {
		if err = batch.Append(dto.GetTransactionFields(&tr)...); err != nil {
			_ = batch.Abort()
			metrics.TransactionErrTotal.Inc()
			return err
		}
	}

	if err = batch.Send(); err != nil {
		_ = batch.Abort()
		metrics.TransactionErrTotal.Add(float64(len(buff)))
		return err
	}
	metrics.TransactionInVec.WithLabelValues("TransactionIn").Add(float64(len(buff)))
	batch = nil

	return nil
}

func (repo *transactionRepo) flushBatch(ctx context.Context, batch []dao.TranProc) []dao.Ss7CdrProc {
	if len(batch) == 0 {
		return batch
	}

	if !repo.GetConnValid() {
		return batch
	}

	batchCtx, cancel := context.WithTimeout(ctx, batchPushTimeout)
	defer cancel()

	if err := repo.runPush(batchCtx, batch); err != nil {
		logger.Error("error pushing transactions: %v", err)
		repo.SetConnValid(false)
		return batch
	}

	return batch[:0]
}

func (repo *transactionRepo) GetConnValid() bool {
	return repo.isConnValid.Load()
}

func (repo *transactionRepo) SetConnValid(valid bool) {
	repo.isConnValid.Store(valid)
}

func (repo *transactionRepo) ConnRecovery(ctx context.Context) error {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if repo.GetConnValid() {
				return nil
			}

			if err := repo.store.Ping(); err == nil {
				repo.SetConnValid(true)
				return nil
			}
		}
	}
}
