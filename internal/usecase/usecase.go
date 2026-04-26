package usecase

import "unit-service/internal/repository"

type Usecase interface {
	GetReferenceUsecase() (ReferenceUsecase, error)
	GetTransactionUsecase() (TransactionUsecase, error)
	GetQueueUsecase() (QueueUsecase, error)
}

type usecase struct {
	repo        repository.Repository
	reference   ReferenceUsecase
	transaction TransactionUsecase
	queue       QueueUsecase
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{repo: repo}
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

	repoQueue, err := u.repo.GetQueue()

	if err != nil {
		return nil, err
	}

	u.queue = NewQueueUsecase(repoQueue, u.transaction)

	return u.queue, nil
}
