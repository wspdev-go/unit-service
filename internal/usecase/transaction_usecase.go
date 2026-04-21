package usecase

import "unit-service/internal/repository"

type TransactionUsecase interface {
	Run() error
}

type transactionUsecase struct {
	repo repository.TransactionRepo
}

func NewTransactionUsecase(repo repository.TransactionRepo) TransactionUsecase {
	return &transactionUsecase{
		repo: repo,
	}
}

func (u *transactionUsecase) Run() error {
	return nil
}
