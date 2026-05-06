package dto

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
