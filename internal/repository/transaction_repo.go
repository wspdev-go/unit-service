package repository

import "unit-service/internal/store"

type TransactionRepo interface {
	PushBatch() error
}

type transactionRepo struct {
	store store.TransactionStore
}

func NewTransactionRepo(store store.TransactionStore) TransactionRepo {
	return &transactionRepo{
		store: store,
	}
}

func (repo *transactionRepo) PushBatch() error {
	return nil
}
