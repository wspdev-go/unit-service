package repository

import "unit-service/internal/store"

type Repository interface {
	GetReference() ReferenceRepo
	GetQueue() QueueRepo
	GetTransaction() TransactionRepo
}

type repo struct {
	store       store.Store
	Reference   ReferenceRepo
	Queue       QueueRepo
	Transaction TransactionRepo
}

func NewRepository(store store.Store) Repository {
	return &repo{
		store: store,
	}
}

func (r *repo) GetReference() ReferenceRepo {
	if r.Reference != nil {
		return r.Reference
	}

	refStore := r.store.GetReference()

	r.Reference = NewReferenceRepo(refStore)

	return r.Reference
}

func (r *repo) GetQueue() QueueRepo {
	if r.Queue != nil {
		return r.Queue
	}

	qStore := r.store.GetQueue()

	r.Queue = NewQueueRepo(qStore)

	return r.Queue
}

func (r *repo) GetTransaction() TransactionRepo {
	if r.Transaction != nil {
		return r.Transaction
	}

	txStore := r.store.GetTransaction()

	r.Transaction = NewTransactionRepo(txStore)

	return r.Transaction
}
