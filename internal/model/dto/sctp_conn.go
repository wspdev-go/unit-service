package dto

type SctpConn struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	LocalInterface     string `json:"local_interface"`
	LocalIpAddress     string `json:"local_ip_address"`
	LocalIpPort        int    `json:"local_ip_port"`
	RemoteIpAddress    string `json:"remote_ip_address"`
	RemoteIpPort       int    `json:"remote_ip_port"`
	SctpRole           string `json:"sctp_role"`
	Heartbeats         bool   `json:"heartbeats"`
	HeartbeatsTimer    uint32 `json:"heartbeats_timer"`
	PathRetransmission uint32 `json:"path_retransmission"`
	MaxAssociations    uint32 `json:"max_associations"`
	NumberOfStreams    uint32 `json:"number_of_streams"`
	IsEnable           bool   `json:"is_enable"`
	Description        string `json:"description"`
	WriteBufferSize    int    `json:"write_buffer_size"`
}

type M3UaAsConn struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	LocalPointCode        string `json:"local_point_code"`
	RemotePointCode       string `json:"remote_point_code"`
	Rc                    int    `json:"rc"`
	NwApr                 int    `json:"nw_apr"`
	Tmt                   int    `json:"tmt"`
	AsType                string `json:"as_type"`
	TrafficMode           string `json:"traffic_mode"`
	SsnmEnabled           int    `json:"ssnm_enabled"`
	IndirectPathDiscovery int    `json:"indirect_path_discovery"`
	IsEnable              bool   `json:"is_enable"`
	Description           string `json:"description"`
}

type M3UaAspLink struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SctpConnID   int    `json:"sctp_conn_id"`
	M3UaAsConnID int    `json:"m3ua_as_conn_id"`
	AspID        int    `json:"asp_id"`
	Sls          int    `json:"sls"`
	AspMode      string `json:"asp_mode"`
	IsEnable     bool   `json:"is_enable"`
	Description  string `json:"description"`
}
