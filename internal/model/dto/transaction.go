package dto

import (
	"time"
	"unit-service/internal/model/dao"
)

type Transaction struct {
	TrID              string `json:"msg_id"`
	TransactionLinkID int    `json:"transaction_link_id"`
	TrDate            string `json:"msg_date"`
}

func ConvertTransaction(tr *Transaction) *dao.Transaction {
	trDate, err := time.Parse("2006-01-02", tr.TrDate)
	if err != nil {
		return nil
	}
	return &dao.Transaction{
		TrID:              tr.TrID,
		TransactionLinkID: tr.TransactionLinkID,
		TrDate:            trDate,
	}
}
