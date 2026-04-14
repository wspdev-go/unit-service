package repository

import (
	"unit-service/internal/model/dao"
	"unit-service/internal/store"
)

type ReferenceRepo interface {
	GetSctpConnList() ([]dao.SctpConn, error)
}

func NewReferenceRepo(store store.ReferenceStore) ReferenceRepo {
	return &referenceRepo{
		store: store,
	}
}

type referenceRepo struct {
	store store.ReferenceStore
}

func (r *referenceRepo) GetSctpConnList() ([]dao.SctpConn, error) {
	sctpConnList := make([]dao.SctpConn, 0)
	return sctpConnList, nil
}
