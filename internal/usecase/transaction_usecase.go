package usecase

import (
	"fmt"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
)

type TransactionUsecase interface {
	Run() error
	Handler(transaction *dto.SS7CDR) error
}

type transactionUsecase struct {
	repo      repository.TransactionRepo
	reference ReferenceReader
}

func NewTransactionUsecase(repo repository.TransactionRepo, refUc ReferenceReader) TransactionUsecase {
	return &transactionUsecase{
		repo:      repo,
		reference: refUc,
	}
}

func (u *transactionUsecase) Run() error {
	// Control ClickHouse connection and prepare transaction processing state here.
	return nil
}

func (u *transactionUsecase) Handler(transaction *dto.SS7CDR) error {
	link, ok := u.reference.GetM3UaLink(transaction.SigtranLinkID)
	if !ok {
		return fmt.Errorf("m3ua link not found: %d", transaction.SigtranLinkID)
	}

	_ = link

	procMess := dto.ConvertSS7CDRToSs7CdrProc(transaction)

	if err := u.repo.PutTransaction(procMess); err != nil {
		return err
	}
	return nil
}
