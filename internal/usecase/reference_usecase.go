/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package usecase

import (
	"context"
	"errors"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
)

type ReferenceReader interface {
	GetM3UaLink(id int) (dto.M3UaLink, bool)
}

type ReferenceUsecase interface {
	Run(ctx context.Context) error
	ReferenceReader
}

type referenceUsecase struct {
	repo            repository.ReferenceRepo
	M3UaLinkList    map[int]dto.M3UaLink
	sctpConnList    map[int]dto.SctpConn
	m3uaAsConnList  map[int]dto.M3UaAsConn
	m3uaAspLinkList map[int]dto.M3UaAspLink
}

func NewReferenceUsecase(repo repository.ReferenceRepo) ReferenceUsecase {
	return &referenceUsecase{
		repo: repo,
	}
}

func (u *referenceUsecase) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := u.GetSctpConnList(); err != nil {
		return err
	} else if len(u.sctpConnList) == 0 {
		return errors.New("sctp conn list is empty")
	}

	if err := u.GetM3UaAsConnList(); err != nil {
		return err
	} else if len(u.m3uaAsConnList) == 0 {
		return errors.New("m3ua_as conn list is empty")
	}

	if err := u.GetM3UaAspLinkList(); err != nil {
		return err
	} else if len(u.m3uaAspLinkList) == 0 {
		return errors.New("m3ua_asp link list is empty")
	}

	u.M3UaLinkList = make(map[int]dto.M3UaLink, len(u.m3uaAspLinkList))

	for _, m3uaAspLink := range u.m3uaAspLinkList {
		var link = dto.M3UaLink{
			SctpConnList:    u.sctpConnList[m3uaAspLink.SctpConnID],
			M3uaAsConnList:  u.m3uaAsConnList[m3uaAspLink.M3UaAsConnID],
			M3uaAspLinkList: m3uaAspLink,
		}
		u.M3UaLinkList[m3uaAspLink.ID] = link

	}

	return nil
}

func (u *referenceUsecase) GetSctpConnList() error {

	sctpConnList, err := u.repo.GetSctpConnList()
	if err != nil {
		return err
	}

	if len(sctpConnList) == 0 {
		return nil
	}

	u.sctpConnList = make(map[int]dto.SctpConn, len(sctpConnList))

	for _, sctpConn := range sctpConnList {
		var conn = dto.SctpConn{
			ID:          int(sctpConn.ID),
			Name:        sctpConn.Name,
			IsEnable:    sctpConn.IsEnable,
			Description: sctpConn.Description,
		}
		u.sctpConnList[conn.ID] = conn
	}

	return nil
}

func (u *referenceUsecase) GetM3UaAsConnList() error {
	m3uaAsConnList, err := u.repo.GetM3uaAsConnList()
	if err != nil {
		return err
	}

	if len(m3uaAsConnList) == 0 {
		return nil
	}

	u.m3uaAsConnList = make(map[int]dto.M3UaAsConn, len(m3uaAsConnList))
	for _, m3uaAsConn := range m3uaAsConnList {
		var conn = dto.M3UaAsConn{
			ID:          int(m3uaAsConn.ID),
			Name:        m3uaAsConn.Name,
			IsEnable:    m3uaAsConn.IsEnable,
			Description: m3uaAsConn.Description,
		}
		u.m3uaAsConnList[conn.ID] = conn
	}

	return nil
}

func (u *referenceUsecase) GetM3UaAspLinkList() error {
	m3uaAspLinkList, err := u.repo.GetM3uaAspLinkList()
	if err != nil {
		return err
	}

	if len(m3uaAspLinkList) == 0 {
		return nil
	}

	u.m3uaAspLinkList = make(map[int]dto.M3UaAspLink, len(m3uaAspLinkList))
	for _, m3uaAspLink := range m3uaAspLinkList {
		var link = dto.M3UaAspLink{
			ID:           int(m3uaAspLink.ID),
			Name:         m3uaAspLink.Name,
			SctpConnID:   int(m3uaAspLink.SctpConnID),
			M3UaAsConnID: int(m3uaAspLink.M3UaAsConnID),
			IsEnable:     m3uaAspLink.IsEnable,
			Description:  m3uaAspLink.Description,
		}
		u.m3uaAspLinkList[link.ID] = link
	}

	return nil
}

func (u *referenceUsecase) GetM3UaLink(id int) (dto.M3UaLink, bool) {
	if u.M3UaLinkList == nil {
		return dto.M3UaLink{}, false
	}

	link, ok := u.M3UaLinkList[id]
	return link, ok
}
