package usecase

import (
	"context"
	"errors"
	"sync"
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

	jobsCh := make(chan *dto.SS7CDR, u.jobQueueSize)
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
		wg.Wait()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		cdr, err := u.repo.Consume(ctx)
		if err != nil {
			return err
		}

		if cdr == nil {
			continue
		}

		select {
		case jobsCh <- cdr:
		case <-ctx.Done():
			return ctx.Err()
		}

		/*if err := u.repo.Publish(ctx, cdr); err != nil {
			logger.Error("failed to publish processed cdr: %v, error: %v", cdr, err)
			continue
		}*/
	}
}

func (u *queueUsecase) runWorker(ctx context.Context, jobsCh <-chan *dto.SS7CDR, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case cdr, ok := <-jobsCh:
			if !ok {
				return
			}

			if err := u.transactionUc.Handler(ctx, cdr); err != nil {
				logger.Error("worker %d failed to process transaction: %v, error: %v", workerID, cdr, err)
				continue
			}

			// later:
			// if err := u.repo.Publish(ctx, cdr); err != nil {
			// 	logger.Error("worker %d failed to publish processed cdr: %v, error: %v", workerID, cdr, err)
			// }
		}
	}
}
