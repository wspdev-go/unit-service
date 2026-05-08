/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package dao

import (
	"time"
)

type Transaction struct {
	TransactionLinkID int       `json:"transaction_link_id"` // Transaction Link ID
	TrID              string    `json:"tr_id"`               // Transaction ID
	TrDate            time.Time `json:"tr_date"`             // Date
}
