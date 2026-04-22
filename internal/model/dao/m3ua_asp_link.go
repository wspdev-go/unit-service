package dao

import "time"

type M3UaAspLink struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	SctpConnID   int64     `gorm:"column:sctp_conn_id" json:"sctp_conn_id"`
	M3UaAsConnID int64     `gorm:"column:m3ua_as_conn_id" json:"m3ua_as_conn_id"`
	AspID        int       `gorm:"column:asp_id" json:"asp_id"`
	Sls          int       `gorm:"column:sls" json:"sls"`
	AspMode      string    `gorm:"column:asp_mode" json:"asp_mode"`
	IsEnable     bool      `gorm:"column:is_enable" json:"is_enable"`
	Description  string    `gorm:"column:description" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
