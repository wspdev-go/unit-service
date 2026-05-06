package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unit-service/internal/model/dto"
	"unit-service/internal/store"

	"github.com/redis/go-redis/v9"
)

const (
	consumerQueueList = "default_cdrfeed" //"consumer_queue_list"
	producerQueueList = "producer_queue_list"
)

type QueueRepo interface {
	Consume(ctx context.Context) (*dto.SS7CDR, error)
	Publish(ctx context.Context, cdr *dto.SS7CDR) error
}

type queueRepo struct {
	client *redis.Client
}

func NewQueueRepo(store store.QueueStore) (QueueRepo, error) {
	redisClient := store.Client()

	if redisClient == nil {
		return nil, errors.New("failed to initialize Redis client")
	}

	return &queueRepo{
		client: redisClient,
	}, nil
}

func (repo *queueRepo) Consume(ctx context.Context) (*dto.SS7CDR, error) {
	if repo.client == nil {
		return nil, errors.New("redis client is nil")
	}

	result, err := repo.client.BRPop(ctx, 3*time.Second, consumerQueueList).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	if len(result) != 2 {
		return nil, fmt.Errorf("unexpected BRPOP result: %v", result)
	}

	var cdr dto.SS7CDR
	if err := json.Unmarshal([]byte(result[1]), &cdr); err != nil {
		return nil, err
	}

	return &cdr, nil
}

func (repo *queueRepo) Publish(ctx context.Context, cdr *dto.SS7CDR) error {
	if repo.client == nil {
		return errors.New("redis client is nil")
	}

	if cdr == nil {
		return errors.New("cdr is nil")
	}

	payload, err := json.Marshal(cdr)
	if err != nil {
		return err
	}

	if err := repo.client.LPush(ctx, producerQueueList, payload).Err(); err != nil {
		return err
	}

	return nil
}
