package usecase

import "unit-service/internal/repository"

type Usecase interface {
	GetReferenceUsecase() ReferenceUsecase
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

func (u *usecase) GetReferenceUsecase() ReferenceUsecase {
	if u.reference != nil {
		return u.reference
	}

	u.reference = NewReferenceUsecase(u.repo.GetReference())

	return u.reference

}

func (u *usecase) GetTransactionUsecase() TransactionUsecase {
	if u.transaction != nil {
		return u.transaction
	}

	u.transaction = NewTransactionUsecase(u.repo.GetTransaction(), u.reference.GetReferenceData())

	return u.transaction
}

func (u *usecase) GetQueueUsecase() QueueUsecase {
	if u.queue != nil {
		return u.queue
	}

	u.queue = NewQueueUsecase(u.repo.GetQueue())

	return u.queue
}
