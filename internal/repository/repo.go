package repository

import "unit-service/internal/store"

type Repository interface {
	GetReference() (ReferenceRepo, error)
	GetQueue() (QueueRepo, error)
	GetTransaction() (TransactionRepo, error)
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

func (r *repo) GetReference() (ReferenceRepo, error) {
	if r.Reference != nil {
		return r.Reference, nil
	}

	var err error

	refStore := r.store.GetReference()

	r.Reference, err = NewReferenceRepo(refStore)
	if err != nil {
		return nil, err
	}

	return r.Reference, nil
}

func (r *repo) GetQueue() (QueueRepo, error) {
	if r.Queue != nil {
		return r.Queue, nil
	}

	var err error

	qStore := r.store.GetQueue()

	r.Queue, err = NewQueueRepo(qStore)

	if err != nil {
		return nil, err
	}

	return r.Queue, nil
}

func (r *repo) GetTransaction() (TransactionRepo, error) {
	if r.Transaction != nil {
		return r.Transaction, nil
	}

	var err error

	txStore := r.store.GetTransaction()

	r.Transaction, err = NewTransactionRepo(txStore)
	if err != nil {
		return nil, err
	}

	return r.Transaction, nil
}
