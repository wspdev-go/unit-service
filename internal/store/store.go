package store

import "unit-service/internal/config"

type Store interface {
	GetReference() ReferenceStore
	GetQueue() QueueStore
	GetTransaction() TransactionStore
}

type store struct {
	cfg         *config.Config
	Reference   ReferenceStore
	Queue       QueueStore
	Transaction TransactionStore
}

func NewStore(cfg *config.Config) Store {
	return store{
		cfg: cfg,
	}
}

func (s store) GetReference() ReferenceStore {
	if s.Reference != nil {
		return s.Reference
	}

	s.Reference = NewReference(s.cfg.ReferenceDB)

	return s.Reference
}

func (s store) GetQueue() QueueStore {
	if s.Queue != nil {
		return s.Queue
	}

	s.Queue = NewQueue(s.cfg.QueueDB)

	return s.Queue
}

func (s store) GetTransaction() TransactionStore {
	if s.Transaction != nil {
		s.Transaction = nil
	}

	s.Transaction = NewTransaction(s.cfg.TransactionDB)

	return s.Transaction
}
