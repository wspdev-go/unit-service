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
