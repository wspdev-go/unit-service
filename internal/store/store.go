package store

import "unit-service/internal/config"

type Store interface {
	GetReference()
	GetQueue()
	GetTransaction()
}

type store struct {
	cfg *config.Config
}

func NewStore(cfg *config.Config) Store {
	return store{
		cfg: cfg,
	}
}

func (s store) GetReference() {}

func (s store) GetQueue() {}

func (s store) GetTransaction() {}
