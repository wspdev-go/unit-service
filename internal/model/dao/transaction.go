package dao

import (
	"time"
)

type Transaction struct {
	TrDate time.Time `json:"tr_date"` // Date
}
