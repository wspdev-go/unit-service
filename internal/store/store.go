package store

import "unit-service/internal/config"

type Store interface {
	GetReference() Reference
	GetQueue() Queue
	GetTransaction() Transaction
}

type store struct {
	cfg         *config.Config
	Reference   Reference
	Queue       Queue
	Transaction Transaction
}

func NewStore(cfg *config.Config) Store {
	return store{
		cfg: cfg,
	}
}

func (s store) GetReference() Reference {
	if s.Reference != nil {
		return s.Reference
	}

	s.Reference = NewReference(s.cfg)

	return s.Reference
}

func (s store) GetQueue() Queue {
	if s.Queue != nil {
		return s.Queue
	}

	s.Queue = NewQueue(s.cfg)

	return s.Queue
}

func (s store) GetTransaction() Transaction {
	if s.Transaction != nil {
		s.Transaction = nil
	}

	s.Transaction = NewTransaction(s.cfg)

	return s.Transaction
}
