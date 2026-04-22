package usecase

import (
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
)

type TransactionUsecase interface {
	Run() error
}

type transactionUsecase struct {
	repo         repository.TransactionRepo
	M3uaLinkList map[int]dto.M3UaLink
}

func NewTransactionUsecase(repo repository.TransactionRepo, links map[int]dto.M3UaLink) TransactionUsecase {
	return &transactionUsecase{
		repo:         repo,
		M3uaLinkList: links,
	}
}

func (u *transactionUsecase) Run() error {
	return nil
}
