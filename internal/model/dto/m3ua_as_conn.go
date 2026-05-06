package dto

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
