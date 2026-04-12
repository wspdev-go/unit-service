package store

import "unit-service/internal/config"

type Reference interface {
	IsOpen() bool
	Close() error
}

type reference struct {
	cfg *config.ReferenceConfig
}

func NewReference(cfg *config.ReferenceConfig) Reference {
	return reference{
		cfg: cfg,
	}
}

func (r reference) IsOpen() bool {
	return r.cfg != nil
}

func (r reference) Close() error {
	return nil
}
