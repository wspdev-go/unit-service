package dto

type M3UaAspLink struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SctpConnID   int    `json:"sctp_conn_id"`
	M3UaAsConnID int    `json:"m3ua_as_conn_id"`
	IsEnable     bool   `json:"is_enable"`
	Description  string `json:"description"`
}
