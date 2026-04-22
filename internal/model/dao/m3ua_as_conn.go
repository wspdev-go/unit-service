package dao

import "time"

type M3UaAsConn struct {
	ID                    int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                  string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	LocalPointCode        string    `gorm:"column:local_point_code" json:"local_point_code"`
	RemotePointCode       string    `gorm:"column:remote_point_code" json:"remote_point_code"`
	Rc                    int       `gorm:"column:rc" json:"rc"`
	NwApr                 int       `gorm:"column:nw_apr" json:"nw_apr"`
	Tmt                   int       `gorm:"column:tmt" json:"tmt"`
	AsType                string    `gorm:"column:as_type" json:"as_type"`
	TrafficMode           string    `gorm:"column:traffic_mode" json:"traffic_mode"`
	SsnmEnabled           int       `gorm:"column:ssnm_enabled" json:"ssnm_enabled"`
	IndirectPathDiscovery int       `gorm:"column:indirect_path_discovery" json:"indirect_path_discovery"`
	IsEnable              bool      `gorm:"column:is_enable" json:"is_enable"`
	Description           string    `gorm:"column:description" json:"description"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
