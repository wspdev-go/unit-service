package usecase

import "unit-service/internal/repository"

type Usecase interface {
	GetReferenceUsecase() ReferenceUsecase
}

type usecase struct {
	repo      repository.Repository
	reference ReferenceUsecase
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
