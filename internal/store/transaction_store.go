package store

import "unit-service/internal/config"

type Transaction interface {
	IsOpen() bool
	Close() error
}

type transaction struct {
	cfg *config.Config
}

func NewTransaction(cfg *config.Config) Transaction {
	return transaction{
		cfg: cfg,
	}
}

func (t transaction) IsOpen() bool {
	return t.cfg.TransactionDB != nil
}

func (t transaction) Close() error {
	return nil
}
