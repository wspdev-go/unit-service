package dao

import (
	"time"
)

type Transaction struct {
	TransactionLinkID int       `json:"transaction_link_id"` // Transaction Link ID
	TrID              string    `json:"tr_id"`               // Transaction ID
	TrDate            time.Time `json:"tr_date"`             // Date
}
