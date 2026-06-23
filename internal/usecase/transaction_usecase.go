/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package usecase

import (
	"context"
	"fmt"
	"time"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
)

type TransactionConnUsecase interface {
	GetConnValid() bool
	ConnRecovery(ctx context.Context) error
}

type TransactionUsecase interface {
	Run(ctx context.Context) error
	Handler(ctx context.Context, transaction *dto.Transaction) error
	TransactionConnUsecase
}

type transactionUsecase struct {
	repo      repository.TransactionRepo
	reference ReferenceReader
}

func NewTransactionUsecase(repo repository.TransactionRepo, refUc ReferenceReader) TransactionUsecase {
	return &transactionUsecase{
		repo:      repo,
		reference: refUc,
	}
}

func (u *transactionUsecase) Run(ctx context.Context) error {
	// Control ClickHouse batch insert by time or count, for example, every 5 seconds or every 100 transactions
	return u.repo.RunBatchWriter(ctx)
}

func (u *transactionUsecase) Handler(ctx context.Context, transaction *dto.Transaction) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	link, ok := u.reference.GetM3UaLink(transaction.TransactionLinkID)
	if !ok {
		return fmt.Errorf("signal link not found: %d", transaction.TransactionLinkID)
	}

	_ = link

	procMess := dto.ConvertTransaction(transaction)

	if err := u.repo.PutBatch(ctx, procMess); err != nil {
		return err
	}
	return nil
}

func (u *transactionUsecase) GetConnValid() bool {
	return u.repo.GetConnValid()
}

func (u *transactionUsecase) ConnRecovery(ctx context.Context) error {
	return u.repo.ConnRecovery(ctx)
}
