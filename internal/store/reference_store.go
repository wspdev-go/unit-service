package store

import "unit-service/internal/config"

type Reference interface {
	IsOpen() bool
	Close() error
}

type reference struct {
	cfg *config.Config
}

func NewReference(cfg *config.Config) Reference {
	return reference{
		cfg: cfg,
	}
}

func (r reference) IsOpen() bool {
	return r.cfg.ReferenceDB != nil
}

func (r reference) Close() error {
	return nil
}
