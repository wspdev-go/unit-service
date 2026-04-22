package usecase

import (
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
)

type TransactionUsecase interface {
	Run() error
	Handler(transaction *dto.SS7CDR) error
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
	//Control ClickHouse Connect and additionally check in CDR
	return nil
}

func (u *transactionUsecase) Handler(transaction *dto.SS7CDR) error {
	procMess := dto.ConvertSS7CDRToSs7CdrProc(transaction)
	if err := u.repo.PutTransaction(procMess); err != nil {
		return err
	}
	return nil
}
