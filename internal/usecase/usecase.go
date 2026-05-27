package usecase

import (
	"context"
	"unit-service/internal/repository"
	"unit-service/internal/store"

	"golang.org/x/sync/errgroup"
)

type Usecase interface {
	RunQueuePipeline(ctx context.Context) error
}

type usecase struct {
	repo        repository.Repository
	reference   ReferenceUsecase
	transaction TransactionUsecase
	queue       QueueUsecase
}

func NewUsecase(store store.Store) Usecase {
	repo := repository.NewRepository(store)
	return &usecase{repo: repo}
}

func (u *usecase) RunQueuePipeline(ctx context.Context) error {
	referenceUc, err := u.GetReferenceUsecase()
	if err != nil {
		return err
	}

	transactionUc, err := u.GetTransactionUsecase()
	if err != nil {
		return err
	}

	queueUc, err := u.GetQueueUsecase()
	if err != nil {
		return err
	}

	// Bootstrap the reference snapshot before starting long-running workers.
	if err := referenceUc.Run(ctx); err != nil {
		return err
	}

	g, groupCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return transactionUc.Run(groupCtx)
	})

	g.Go(func() error {
		return queueUc.Run(groupCtx)
	})

	return g.Wait()
}

func (u *usecase) GetReferenceUsecase() (ReferenceUsecase, error) {
	if u.reference != nil {
		return u.reference, nil
	}

	refRepo, err := u.repo.GetReference()
	if err != nil {
		return nil, err
	}

	u.reference = NewReferenceUsecase(refRepo)

	return u.reference, nil

}

func (u *usecase) GetTransactionUsecase() (TransactionUsecase, error) {
	if u.transaction != nil {
		return u.transaction, nil
	}

	var err error

	u.reference, err = u.GetReferenceUsecase()
	if err != nil {
		return nil, err
	}

	repoTransaction, err := u.repo.GetTransaction()
	if err != nil {
		return nil, err
	}

	u.transaction = NewTransactionUsecase(repoTransaction, u.reference)

	return u.transaction, nil
}

func (u *usecase) GetQueueUsecase() (QueueUsecase, error) {
	if u.queue != nil {
		return u.queue, nil
	}

	var err error

	u.transaction, err = u.GetTransactionUsecase()
	if err != nil {
		return nil, err
	}

	repoQueue, err := u.repo.GetQueue()

	if err != nil {
		return nil, err
	}

	u.queue = NewQueueUsecase(repoQueue, u.transaction)

	return u.queue, nil
}
