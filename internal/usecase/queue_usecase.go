package usecase

import "unit-service/internal/repository"

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
	return nil
}
