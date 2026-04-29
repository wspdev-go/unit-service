package usecase

import (
	"context"
	"errors"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
	"unit-service/logger"
)

type QueueUsecase interface {
	Run(ctx context.Context) error
}

type queueUsecase struct {
	repo          repository.QueueRepo
	transactionUc TransactionUsecase
}

func NewQueueUsecase(repo repository.QueueRepo, trUc TransactionUsecase) QueueUsecase {
	return &queueUsecase{
		repo:          repo,
		transactionUc: trUc,
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

	semCh := make(chan struct{}, 10)

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

		semCh <- struct{}{}

		go func(cdr *dto.SS7CDR) {
			defer func() { <-semCh }()
			if err = u.transactionUc.Handler(ctx, cdr); err != nil {
				logger.Error("failed to process transaction: %v, error: %v", cdr, err)
			}
		}(cdr)

		/*if err := u.repo.Publish(ctx, cdr); err != nil {
			logger.Error("failed to publish processed cdr: %v, error: %v", cdr, err)
			continue
		}*/
	}
}
