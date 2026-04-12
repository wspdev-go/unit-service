package store

import "unit-service/internal/config"

type Queue interface {
	IsOpen() bool
	Close() error
}

type queue struct {
	cfg *config.QueueConfig
}

func NewQueue(cfg *config.QueueConfig) Queue {
	return queue{
		cfg: cfg,
	}
}

func (q queue) IsOpen() bool {
	return q.cfg != nil
}

func (q queue) Close() error {
	return nil
}
