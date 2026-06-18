package usecase

import (
	"context"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
)

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
	return u.repo.RunBatchWriter(ctx)
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

	if err := u.repo.PutBatch(ctx, procMess); err != nil {
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
