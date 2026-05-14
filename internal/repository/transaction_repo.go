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
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unit-service/internal/model/dao"
	"unit-service/internal/store"

	"github.com/ClickHouse/clickhouse-go/v2"
)

const batchSize = 1000

type TransactionRepo interface {
	PushBatchTransaction(ctx context.Context) error
	PutBatch(transaction *dao.Transaction) error
	PutTransaction(transaction *dao.Transaction) error
	FlushCh() <-chan struct{}
	GetConnValid() bool
	SetConnValid(valid bool)
	ConnRecovery(ctx context.Context) error
}

type transactionRepo struct {
	conn        clickhouse.Conn
	store       store.TransactionStore
	batchBuff   []dao.Transaction
	mu          sync.Mutex
	flushCh     chan struct{}
	isConnValid atomic.Bool
}

func NewTransactionRepo(store store.TransactionStore) (TransactionRepo, error) {
	conn := store.Conn()

	if conn == nil {
		return nil, errors.New("clickhouse connection is nil")
	}

	return &transactionRepo{
		conn:      conn,
		store:     store,
		batchBuff: make([]dao.Transaction, 0, batchSize), // Initialize batch buffer with a reasonable capacity
		flushCh:   make(chan struct{}, 1),
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

func gettrFields(tr *dao.Transaction) []any {
	return []any{
		tr.TrDate,
	}
}

func getRepoInsQuery(obj dao.Transaction) string {
	t := reflect.TypeOf(obj)

	columns := make([]string, 0)
	placeholders := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get("json")
		if col == "" || col == "-" {
			continue
		}

		columns = append(columns, strings.ReplaceAll(col, ",omitempty", ""))
		placeholders = append(placeholders, "?")
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		" tr",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return sql
}

func (repo *transactionRepo) PutBatch(transaction *dao.Transaction) error {
	if transaction == nil {
		return errors.New("transaction is nil")
	}

	repo.mu.Lock()
	repo.batchBuff = append(repo.batchBuff, *transaction)
	if len(repo.batchBuff) >= batchSize {
		select {
		case repo.flushCh <- struct{}{}:
		default:
		}
	}
	repo.mu.Unlock()

	return nil
}

func (repo *transactionRepo) FlushCh() <-chan struct{} {
	return repo.flushCh
}

func (repo *transactionRepo) PushBatchTransaction(ctx context.Context) error {
	repo.mu.Lock()

	if len(repo.batchBuff) == 0 {
		repo.mu.Unlock()
		return nil
	}

	if !repo.GetConnValid() {
		repo.mu.Unlock()
		return nil
	}

	buff := make([]dao.Transaction, len(repo.batchBuff))
	copy(buff, repo.batchBuff)
	repo.batchBuff = make([]dao.Transaction, 0, 3*batchSize)

	repo.mu.Unlock()

	return repo.runPush(ctx, buff)
}

func (repo *transactionRepo) runPush(ctx context.Context, buff []dao.Transaction) error {
	batch, err := repo.conn.PrepareBatch(ctx, getRepoInsQuery(dao.Transaction{}))
	if err != nil {
		repo.restoreFailedBatch(buff)
		return err
	}

	for _, tr := range buff {
		if err = batch.Append(gettrFields(&tr)...); err != nil {
			_ = batch.Abort()
			repo.restoreFailedBatch(buff)
			return err
		}
	}

	if err = batch.Send(); err != nil {
		_ = batch.Abort()
		repo.restoreFailedBatch(buff)
		return err
	}

	batch = nil

	return nil
}

func (repo *transactionRepo) restoreFailedBatch(buff []dao.Transaction) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.batchBuff = append(buff, repo.batchBuff...)
}

func (repo *transactionRepo) GetConnValid() bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	return repo.isConnValid
}

func (repo *transactionRepo) SetConnValid(valid bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.isConnValid = valid
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
