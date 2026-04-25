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
	repo        repository.TransactionRepo
	referenceUc ReferenceUsecase
}

func NewTransactionUsecase(repo repository.TransactionRepo, refUc ReferenceUsecase) TransactionUsecase {
	return &transactionUsecase{
		repo:        repo,
		referenceUc: refUc,
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
