package repository

import (
	"unit-service/internal/model/dao"
	"unit-service/internal/store"
)

type ReferenceRepo interface {
	GetSctpConnList() ([]dao.SctpConn, error)
	GetM3uaAsConnList() ([]dao.M3UaAsConn, error)
	GetM3uaAspLinkList() ([]dao.M3UaAspLink, error)
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

	db := r.store.DB()

	if db == nil {
		return sctpConnList, nil
	}

	if err := db.Where("is_enable = 1").Find(&sctpConnList).Error; err != nil {
		return nil, err
	}

	return sctpConnList, nil
}

func (r *referenceRepo) GetM3uaAsConnList() ([]dao.M3UaAsConn, error) {
	m3uaAsConnList := make([]dao.M3UaAsConn, 0)

	db := r.store.DB()

	if db == nil {
		return m3uaAsConnList, nil
	}

	if err := db.Table("m3ua_as_conns").Where("is_enable = 1").Find(&m3uaAsConnList).Error; err != nil {
		return nil, err
	}

	return m3uaAsConnList, nil
}

func (r *referenceRepo) GetM3uaAspLinkList() ([]dao.M3UaAspLink, error) {
	m3uaAspLinkList := make([]dao.M3UaAspLink, 0)

	db := r.store.DB()
	if db == nil {
		return m3uaAspLinkList, nil
	}
	if err := db.Table("m3ua_asp_links").Where("is_enable = 1").Find(&m3uaAspLinkList).Error; err != nil {
		return nil, err
	}

	return m3uaAspLinkList, nil
}
