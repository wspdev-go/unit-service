package usecase

import (
	"time"
	"unit-service/internal/repository"
	"unit-service/logger"
)

type QueueUsecase interface {
	Run() error
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

func (u *queueUsecase) Run() error {
	// Run worker that will read from queue and process transactions
	for {
		cdrs, err := u.repo.Get()
		if err != nil {
			return err
		}

		if len(cdrs) == 0 {
			time.Sleep(3 * time.Second)
		}

		for _, cdr := range cdrs {
			err = u.transactionUc.Handler(&cdr)
			if err != nil {
				logger.Error("Failed to process transaction: %v, error: %v", cdr, err)
			}
		}
	}
}
