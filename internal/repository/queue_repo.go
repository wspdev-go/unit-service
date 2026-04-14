package repository

import "unit-service/internal/store"

type QueueRepo interface {
	Put() error
	Get() error
}

type queueRepo struct {
	store store.QueueStore
}

func NewQueueRepo(store store.QueueStore) QueueRepo {
	return &queueRepo{
		store: store,
	}
}

func (repo *queueRepo) Put() error {
	return nil
}

func (repo *queueRepo) Get() error {
	return nil
}
