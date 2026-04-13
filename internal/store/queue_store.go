package store

import (
	"context"
	"fmt"
	"time"
	"unit-service/internal/config"
	"unit-service/logger"

	"github.com/redis/go-redis/v9"
)

type QueueStore interface {
	Open() bool
	Close() error
	Ping() error
}

type queue struct {
	cfg    *config.QueueConfig
	Client *redis.Client
}

func NewQueue(cfg *config.QueueConfig) QueueStore {
	return &queue{
		cfg: cfg,
	}
}

func (q *queue) Open() bool {
	if q.Client != nil {
		return true
	}

	if q.cfg == nil {
		logger.Error("queue config is nil")
		return false
	}

	if q.cfg.Host == "" {
		return false
	}

	if q.cfg.Port == 0 {
		return false
	}

	redisAddr := fmt.Sprintf("%s:%d", q.cfg.Host, q.cfg.Port)

	q.Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: q.cfg.Username,
		Password: q.cfg.Password,
		DB:       q.cfg.Database,
		PoolSize: q.cfg.PoolSize,
	})

	if err := q.Ping(); err != nil {
		logger.Error("failed to connect to queue: %s", err.Error())
		_ = q.Close()
		return false
	}

	return true
}

func (q *queue) Ping() error {
	if q.Client == nil {
		return fmt.Errorf("queue store is not open")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := q.Client.Ping(ctx).Result()
	return err
}

func (q *queue) Close() error {
	if q.Client == nil {
		return nil
	}

	if err := q.Client.Close(); err != nil {
		logger.Error("Failed to close Transaction connection: %s", err.Error())
		return err
	}

	logger.Info("close queue successfully")
	return nil
}
