package repository

import (
	"context"
	"encoding/json"
	"errors"
	"unit-service/internal/model/dto"
	"unit-service/internal/store"

	"unit-service/logger"

	"github.com/redis/go-redis/v9"
)

const (
	consumerQueueList = "default_cdrfeed" //"consumer_queue_list"
	producerQueueList = "producer_queue_list"
)

type QueueRepo interface {
	Put() error
	Get() ([]dto.SS7CDR, error)
}

type queueRepo struct {
	store store.QueueStore
}

func NewQueueRepo(store store.QueueStore) QueueRepo {
	return &queueRepo{
		store: store,
	}
}

func (repo *queueRepo) Put() error {
	return nil
}

func (repo *queueRepo) Get() ([]dto.SS7CDR, error) {
	client := repo.store.Client()
	if client == nil {
		return nil, errors.New("client is nil")
	}

	var cdrs []dto.SS7CDR
	cdrList, err := client.LRange(context.Background(), consumerQueueList, 0, 100).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	for _, val := range cdrList {
		var cdr dto.SS7CDR
		if err = json.Unmarshal([]byte(val), &cdr); err != nil {
			logger.Error("Failed to unmarshal CDR: %s", val)
			return nil, err
		}
		cdrs = append(cdrs, cdr)
		if _, err = client.RPop(context.Background(), consumerQueueList).Result(); err != nil {
			return nil, err
		}
	}

	return cdrs, nil
}
