package dao

import "time"

type M3UaAsConn struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	IsEnable    bool      `gorm:"column:is_enable" json:"is_enable"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
