package dao

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
)

type Ss7CdrProc struct {
	// Date fields
	MsgDate     time.Time `json:"msg_date"`                // Date
	MsgDtUs     time.Time `json:"msg_dt_us"`               // DateTime64(6, 'UTC')
	MsgExpiryDt time.Time `json:"msg_expiry_dt,omitempty"` // DateTime('UTC')

	// Message IDs
	ExtMsgID      string `json:"ext_msg_id"`      // String
	ProxyMsgID    string `json:"proxy_msg_id"`    // String
	InternalMsgID string `json:"internal_msg_id"` // String
	TranMsgID     string `json:"tran_msg_id"`     // String

	// IP and Port information
	SrcIP   string `json:"src_ip"`   // IPv4
	SrcPort uint16 `json:"src_port"` // UInt16
	DstIP   string `json:"dst_ip"`   // IPv4
	DstPort uint16 `json:"dst_port"` // UInt16

	// Message types and kinds
	MsgType int8 `json:"msg_type"` // Int8 DEFAULT -1
	MsgKind int8 `json:"msg_kind"` // Int8 DEFAULT -1
	MsuType int8 `json:"msu_type"` // Int8 DEFAULT -1
	Type    int8 `json:"type"`     // Int8 DEFAULT -1

	// Direction
	Direction int8 `json:"direction"` // Int8 DEFAULT -1

	// Result information
	ResultCode   int16  `json:"result_code,omitempty"` // Int16 DEFAULT -1
	ResultStatus string `json:"result_status"`         // LowCardinality(String)

	// Message addresses
	SenderOA      string `json:"sender_oa"`      // String
	DestinationDA string `json:"destination_da"` // String

	OPC uint32 `json:"opc"` // UInt32
	DPC uint32 `json:"dpc"` // UInt32

	SccpCarrier  string `json:"sccp_carrier"`            // LowCardinality(String)
	SccpClgpaGt  string `json:"sccp_clgpa_gt"`           // String
	SccpClgpaTt  int16  `json:"sccp_clgpa_tt,omitempty"` // Int16 DEFAULT -1
	SccpClgpaSsn int16  `json:"sccp_clgpa_ssn"`          // Int16 DEFAULT -1
	SccpCldpaGt  string `json:"sccp_cldpa_gt"`           // String
	SccpCldpaTt  int16  `json:"sccp_cldpa_tt,omitempty"` // Int16 DEFAULT -1
	SccpCldpaSsn int16  `json:"sccp_cldpa_ssn"`          // Int16 DEFAULT -1

	TcapID uint32 `json:"tcap_id,omitempty"` // UInt32

	MapScentreAddr *big.Int `json:"map_scentre_addr"`   // String containing UInt128 as digits only
	MapMscGt       *big.Int `json:"map_msc_gt"`         // String containing UInt128 as digits only
	MapImsi        string   `json:"map_imsi"`           // String
	MapOpco        int16    `json:"map_opco,omitempty"` // Int16 DEFAULT -1

	CustomerAccount   string `json:"customer_account"`    // LowCardinality(String)
	CustomerAccountID uint32 `json:"customer_account_id"` // UInt32
	SupplierAccount   string `json:"supplier_account"`    // LowCardinality(String)
	SupplierAccountID uint32 `json:"supplier_account_id"` // UInt32 (Note: typo in schema)

	SignallingConnLink   string `json:"signalling_conn_link"`    // LowCardinality(String)
	SignallingConnLinkID uint32 `json:"signalling_conn_link_id"` // UInt32

	DestinationCountry    string `json:"destination_country"`     // LowCardinality(String)
	DestinationCountryID  uint32 `json:"destination_country_id"`  // UInt32
	DestinationOperator   string `json:"destination_operator"`    // String
	DestinationOperatorID uint32 `json:"destination_operator_id"` // UInt32

	EsmClass    int16  `json:"esm_class,omitempty"`     // Int16 DEFAULT -1
	DataCoding  int16  `json:"data_coding,omitempty"`   // Int16 DEFAULT -1
	Pid64       int8   `json:"pid64,omitempty"`         // Int8 DEFAULT -1
	MsgTextLen  int16  `json:"msg_text_len,omitempty"`  // Int16 DEFAULT -1
	Udh         int8   `json:"udh,omitempty"`           // Int8 DEFAULT -1
	MsgRefNum   uint16 `json:"msg_ref_num,omitempty"`   // UInt16
	MsgTotalNum int8   `json:"msg_total_num,omitempty"` // Int8 DEFAULT -1
	MsgPartNum  int8   `json:"msg_part_num,omitempty"`  // Int8 DEFAULT -1

	// DLR information
	DlrErr  int16  `json:"dlr_err,omitempty"` // Int16 DEFAULT -1
	DlrStat string `json:"dlr_stat"`          // LowCardinality(String)

	// Retry information
	RetryPattern int8  `json:"retry_pattern,omitempty"` // Int8 DEFAULT -1
	RetryError   int16 `json:"retry_error,omitempty"`   // Int16 DEFAULT -1
	RetryAttempt int16 `json:"retry_attempt"`           // Int16 DEFAULT -1

	RoutingType          int8   `json:"routing_type"`           // Int8 DEFAULT -1
	TransformationRuleID uint16 `json:"transformation_rule_id"` // UInt16

	MsgData    string `json:"msg_data"`
	MsgDataBin string `json:"msg_data_bin"`
	UdhData    string `json:"udh_data"`
	UdhDataBin string `json:"udh_data_bin"`
}

func GetRepoInsQuery() string {
	t := reflect.TypeOf(Ss7CdrProc{})

	columns := make([]string, 0)
	placeholders := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get("json")
		if col == "" || col == "-" {
			continue
		}

		columns = append(columns, strings.ReplaceAll(col, ",omitempty", ""))
		placeholders = append(placeholders, "?")
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		" cdr",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return sql
}
