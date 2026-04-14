package dao

type SctpConn struct {
	ID                 int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name               string `gorm:"type:varchar(255);not null;unique" json:"name"`
	LocalInterface     string `gorm:"column:local_interface"`
	LocalIpAddress     string `gorm:"column:local_ip_addr"`
	LocalIpPort        int    `gorm:"column:local_port"`
	RemoteIpAddress    string `gorm:"column:remote_ip_addr"`
	RemoteIpPort       int    `gorm:"column:remote_port"`
	SctpRole           string `gorm:"column:sctp_role"`
	Heartbeats         bool   `gorm:"column:heartbeats"`
	HeartbeatsTimer    uint32 `gorm:"column:heartbeats_timer"`
	PathRetransmission uint32 `gorm:"column:path_retransmission"`
	MaxAssociations    uint32 `gorm:"column:max_associations"`
	NumberOfStreams    uint32 `gorm:"column:number_of_streams"`
	IsEnable           bool   `gorm:"column:is_enable"`
	Description        string `gorm:"column:description"`
	WriteBufferSize    int    `gorm:"column:write_buffer_size"`
	CreatedAt          int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

/*
CREATE TABLE `sctp_conns` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `local_interface` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `local_ip_addr` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `local_port` int unsigned NOT NULL DEFAULT '0',
  `remote_ip_addr` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `remote_port` int unsigned NOT NULL DEFAULT '0',
  `sctp_role` enum('server','client') CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT 'server',
  `heartbeats` tinyint(1) DEFAULT '1',
  `heartbeats_timer` int unsigned NOT NULL DEFAULT '0',
  `path_retransmission` int unsigned NOT NULL DEFAULT '0',
  `max_associations` int unsigned NOT NULL DEFAULT '0',
  `number_of_streams` int unsigned NOT NULL DEFAULT '0',
  `is_enable` tinyint(1) DEFAULT '1',
  `description` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_updated_at` (`updated_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3;
*/
