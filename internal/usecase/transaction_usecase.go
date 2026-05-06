package usecase

import (
	"context"
	"time"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
	"unit-service/logger"
)

const batchTimeout = 300

type TransactionConnUsecase interface {
	GetConnValid() bool
	ConnRecovery(ctx context.Context) error
}

type TransactionUsecase interface {
	Run(ctx context.Context) error
	Handler(ctx context.Context, transaction *dto.SS7CDR) error
	TransactionConnUsecase
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
			if err := u.pushTransaction(ctx); err != nil {
				logger.Error("error pushing transactions: %v", err)
				u.repo.SetConnValid(false)
			}
		case <-u.repo.FlushCh():
			if err := u.pushTransaction(ctx); err != nil {
				logger.Error("error pushing transactions: %v", err)
				u.repo.SetConnValid(false)
			}
		}
	}
}

func (u *transactionUsecase) pushTransaction(ctx context.Context) error {
	if !u.repo.GetConnValid() {
		return nil
	}
	batchCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := u.repo.PushBatchTransaction(batchCtx); err != nil {
		return err
	}
	return nil
}

func (u *transactionUsecase) Handler(ctx context.Context, transaction *dto.SS7CDR) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	link, ok := u.reference.GetM3UaLink(transaction.SigtranLinkID)
	if !ok {
		//return fmt.Errorf("Signal link not found: %d", transaction.SigtranLinkID)
	}

	_ = link

	procMess := dto.ConvertSS7CDRToSs7CdrProc(transaction)

	if err := u.repo.PutBatch(procMess); err != nil {
		return err
	}
	return nil
}

func (u *transactionUsecase) GetConnValid() bool {
	return u.repo.GetConnValid()
}

func (u *transactionUsecase) ConnRecovery(ctx context.Context) error {
	return u.repo.ConnRecovery(ctx)
}
