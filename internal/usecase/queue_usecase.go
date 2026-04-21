package usecase

import "unit-service/internal/repository"

type QueueUsecase interface {
	Run() error
}

type queueUsecase struct {
	repo repository.QueueRepo
}

func NewQueueUsecase(repo repository.QueueRepo) QueueUsecase {
	return &queueUsecase{
		repo: repo,
	}
}

func (u *queueUsecase) Run() error {
	return nil
}
