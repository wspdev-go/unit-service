package usecase

import (
	"context"
	"fmt"
	"time"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
	"unit-service/logger"
)

const batchTimeout = 300

type TransactionUsecase interface {
	Run(ctx context.Context) error
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

func (u *transactionUsecase) Run(ctx context.Context) error {
	// Control ClickHouse batch insert by time or count, for example, every 5 seconds or every 100 transactions

	ticker := time.NewTicker(batchTimeout * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			batchCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			err := u.repo.PushBatch(batchCtx)
			if err != nil {
				logger.Error("error pushing transactions: %v", err)
			}
			cancel()
		}
	}
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
