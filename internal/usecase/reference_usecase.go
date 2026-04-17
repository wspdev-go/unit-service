package usecase

import (
	"errors"
	"unit-service/internal/model/dto"
	"unit-service/internal/repository"
)

type ReferenceUsecase interface {
	Run() error
	GetSctpConnList() error
}

type referenceUsecase struct {
	repo         repository.ReferenceRepo
	sctpConnList map[int]dto.SctpConn
}

func NewReferenceUsecase(repo repository.ReferenceRepo) ReferenceUsecase {
	return &referenceUsecase{
		repo: repo,
	}
}

func (u *referenceUsecase) Run() error {
	if err := u.GetSctpConnList(); err != nil {
		return err
	} else if len(u.sctpConnList) == 0 {
		return errors.New("sctp conn list is empty")
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
			ID:                 int(sctpConn.ID),
			Name:               sctpConn.Name,
			LocalInterface:     sctpConn.LocalInterface,
			LocalIpAddress:     sctpConn.LocalIpAddress,
			LocalIpPort:        sctpConn.LocalIpPort,
			RemoteIpAddress:    sctpConn.RemoteIpAddress,
			RemoteIpPort:       sctpConn.RemoteIpPort,
			SctpRole:           sctpConn.SctpRole,
			Heartbeats:         sctpConn.Heartbeats,
			HeartbeatsTimer:    sctpConn.HeartbeatsTimer,
			PathRetransmission: sctpConn.PathRetransmission,
			MaxAssociations:    sctpConn.MaxAssociations,
			NumberOfStreams:    sctpConn.NumberOfStreams,
			IsEnable:           sctpConn.IsEnable,
			Description:        sctpConn.Description,
			WriteBufferSize:    sctpConn.WriteBufferSize,
		}
		u.sctpConnList[conn.ID] = conn
	}

	return nil
}
