/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

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
	return &store{
		cfg: cfg,
	}
}

func (s *store) GetReference() ReferenceStore {
	if s.Reference != nil {
		return s.Reference
	}

	s.Reference = NewReference(s.cfg.ReferenceDB)

	return s.Reference
}

func (s *store) GetQueue() QueueStore {
	if s.Queue != nil {
		return s.Queue
	}

	s.Queue = NewQueue(s.cfg.QueueDB)

	return s.Queue
}

func (s *store) GetTransaction() TransactionStore {
	if s.Transaction != nil {
		return s.Transaction
	}

	s.Transaction = NewTransaction(s.cfg.TransactionDB)

	return s.Transaction
}
