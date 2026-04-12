package store

import "unit-service/internal/config"

type Queue interface {
	IsOpen() bool
	Close() error
}

type queue struct {
	cfg *config.Config
}

func NewQueue(cfg *config.Config) Queue {
	return queue{
		cfg: cfg,
	}
}

func (q queue) IsOpen() bool {
	return q.cfg.QueueDB != nil
}

func (q queue) Close() error {
	return nil
}
