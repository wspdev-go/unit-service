/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package store

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unit-service/internal/config"
	"unit-service/logger"

	"github.com/redis/go-redis/v9"
)

type QueueStore interface {
	Open() error
	Close() error
	Ping() error
	Client() *redis.Client
}

type queue struct {
	cfg    *config.QueueConfig
	client *redis.Client
}

func NewQueue(cfg *config.QueueConfig) QueueStore {
	return &queue{
		cfg: cfg,
	}
}

func (q *queue) Open() error {
	if q.client != nil {
		return nil
	}

	if q.cfg == nil {
		return errors.New("queue config is nil")
	}

	if q.cfg.Host == "" {
		return errors.New("host is empty")
	}

	if q.cfg.Port == 0 {
		return errors.New("port is empty")
	}

	redisAddr := fmt.Sprintf("%s:%d", q.cfg.Host, q.cfg.Port)

	q.client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: q.cfg.Username,
		Password: q.cfg.Password,
		DB:       q.cfg.Database,
		PoolSize: q.cfg.PoolSize,
	})

	if err := q.Ping(); err != nil {
		logger.Error("failed to connect to queue: %s", err.Error())
		_ = q.Close()
		return err
	}

	return nil
}

func (q *queue) Ping() error {
	if q.client == nil {
		return fmt.Errorf("queue store is not open")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := q.client.Ping(ctx).Result()
	return err
}

func (q *queue) Close() error {
	if q.client == nil {
		return nil
	}

	if err := q.client.Close(); err != nil {
		logger.Error("Failed to close Transaction connection: %s", err.Error())
		return err
	}
	q.client = nil

	logger.Info("close queue successfully")
	return nil
}

func (q *queue) Client() *redis.Client {
	return q.client
}
