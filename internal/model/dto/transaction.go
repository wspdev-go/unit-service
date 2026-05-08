/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

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
