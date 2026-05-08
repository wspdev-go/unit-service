/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package usecase

import (
	"context"
	"errors"
	"sync"
	"time"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
	"unit-service/logger"
)

const (
	defaultWorkerCount  = 10
	defaultJobQueueSize = 100
)

type QueueUsecase interface {
	Run(ctx context.Context) error
}

type queueUsecase struct {
	repo          repository.QueueRepo
	transactionUc TransactionUsecase
	workerCount   int
	jobQueueSize  int
}

func NewQueueUsecase(repo repository.QueueRepo, trUc TransactionUsecase) QueueUsecase {
	return &queueUsecase{
		repo:          repo,
		transactionUc: trUc,
		workerCount:   defaultWorkerCount,
		jobQueueSize:  defaultJobQueueSize,
	}
}

func (u *queueUsecase) Run(ctx context.Context) error {
	// Run worker that will read from queue and process transactions
	if u.repo == nil {
		return errors.New("repository is nil")
	}

	if u.transactionUc == nil {
		return errors.New("transactionUc is nil")
	}

	jobsCh := make(chan *dto.Transaction, u.jobQueueSize)
	var wg sync.WaitGroup

	for i := 0; i < u.workerCount; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			u.runWorker(ctx, jobsCh, workerID)
		}(i)
	}

	defer func() {
		close(jobsCh)
	}()

	var err error

loop:
	for {

		if !u.transactionUc.GetConnValid() {
			if err = u.transactionUc.ConnRecovery(ctx); err != nil {
				break loop
			}
		}

		select {
		case <-ctx.Done():
			err = ctx.Err()
			break loop
		default:
		}

		tr, err := u.repo.Consume(ctx)
		if err != nil {
			break loop
		}

		if tr == nil {
			continue
		}

		select {
		case jobsCh <- tr:
		case <-ctx.Done():
			break loop
		}

		/*if err := u.repo.Publish(ctx, tr); err != nil {
			logger.Error("failed to publish processed tr: %v, error: %v", tr, err)
			continue
		}*/
	}
	wg.Wait()

	return err
}

func (u *queueUsecase) runWorker(ctx context.Context, jobsCh <-chan *dto.Transaction, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case tr, ok := <-jobsCh:
			if !ok {
				return
			}
			ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
			err := u.transactionUc.Handler(ctxTimeout, tr)
			if err != nil {
				logger.Error("worker %d failed to process transaction: %v, error: %v", workerID, tr, err)
			}
			cancel()
		}

		// later:
		// if err := u.repo.Publish(ctx, tr); err != nil {
		// 	logger.Error("worker %d failed to publish processed tr: %v, error: %v", workerID, tr, err)
		// }
	}
}
