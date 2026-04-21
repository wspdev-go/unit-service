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
