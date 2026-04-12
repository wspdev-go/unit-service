package store

import "unit-service/internal/config"

type Transaction interface {
	IsOpen() bool
	Close() error
}

type transaction struct {
	cfg *config.TransactionConfig
}

func NewTransaction(cfg *config.TransactionConfig) Transaction {
	return transaction{
		cfg: cfg,
	}
}

func (t transaction) IsOpen() bool {
	return t.cfg != nil
}

func (t transaction) Close() error {
	return nil
}
