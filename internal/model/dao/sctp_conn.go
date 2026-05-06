package dao

import "time"

type SctpConn struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name               string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	LocalInterface     string    `gorm:"column:local_interface"`
	LocalIpAddress     string    `gorm:"column:local_ip_addr"`
	LocalIpPort        int       `gorm:"column:local_port"`
	RemoteIpAddress    string    `gorm:"column:remote_ip_addr"`
	RemoteIpPort       int       `gorm:"column:remote_port"`
	SctpRole           string    `gorm:"column:sctp_role"`
	Heartbeats         bool      `gorm:"column:heartbeats"`
	HeartbeatsTimer    uint32    `gorm:"column:heartbeats_timer"`
	PathRetransmission uint32    `gorm:"column:path_retransmission"`
	MaxAssociations    uint32    `gorm:"column:max_associations"`
	NumberOfStreams    uint32    `gorm:"column:number_of_streams"`
	IsEnable           bool      `gorm:"column:is_enable"`
	Description        string    `gorm:"column:description"`
	WriteBufferSize    int       `gorm:"column:write_buffer_size"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
