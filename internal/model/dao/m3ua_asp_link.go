/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package dao

import "time"

type M3UaAspLink struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	SctpConnID   int64     `gorm:"column:sctp_conn_id" json:"sctp_conn_id"`
	M3UaAsConnID int64     `gorm:"column:m3ua_as_conn_id" json:"m3ua_as_conn_id"`
	IsEnable     bool      `gorm:"column:is_enable" json:"is_enable"`
	Description  string    `gorm:"column:description" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
