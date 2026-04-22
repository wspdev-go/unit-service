package dto

import (
	"math/big"
	"net"
	"time"
	"unit-service/internal/model/dao"
)

type CDRMessageType int

const (
	CDRMessageUndefined CDRMessageType = iota - 1
	_

	CDRMessageHttpSubmitSm
	CDRMessageHttpDeliverSm

	CDRMessageSmppSubmitSm
	CDRMessageSmppDeliverSm

	CDRMessageHttpMnpSriSm

	CDRMessageSs7SriSm
	CDRMessageSs7FsmMt
	CDRMessageSs7FsmMo

	CDRMessageSs7MnpSriSm

	CDRMessageSs7SriSmSpoofCheck
	SendRoutingInfoForSMOpCode = 0x2D
	MOForwardSMOpCode          = 0x2E
)

var messageType = map[CDRMessageType]string{
	CDRMessageUndefined: "undefined",

	CDRMessageSmppSubmitSm:  "SMPP_Submit_SM",
	CDRMessageSmppDeliverSm: "SMPP_Deliver_SM",

	CDRMessageHttpMnpSriSm: "HTTP_MNP_SRI",

	CDRMessageSs7SriSm: "SRI_SM",
	CDRMessageSs7FsmMt: "FSM_MT",
	CDRMessageSs7FsmMo: "FSM_MO",

	CDRMessageSs7MnpSriSm: "SS7_MNP_SRI",

	CDRMessageSs7SriSmSpoofCheck: "SRI_SM_Spoof_Check",
}

type CDRMessageKind int

const (
	CDRMessageKindUndefined CDRMessageKind = iota - 1
	_
	CDRMessageKindRequest
	CDRMessageKindResponse
)

var messageKind = map[CDRMessageKind]string{
	CDRMessageKindUndefined: "undefined",
	CDRMessageKindRequest:   "Request",
	CDRMessageKindResponse:  "Response",
}

type CDRTrafficDirection int

const (
	CDRTrafficDirectionUndefined CDRTrafficDirection = iota - 1
	_
	CDRTrafficDirectionInbound
	CDRTrafficDirectionOutbound
)

var trafficDirection = map[CDRTrafficDirection]string{
	CDRTrafficDirectionUndefined: "undefined",
	CDRTrafficDirectionInbound:   "Inbound",
	CDRTrafficDirectionOutbound:  "Outbound",
}

type CDRTrafficType int

const (
	CDRTrafficTypeUndefined CDRTrafficType = iota - 1
	_
	CDRTrafficTypeSMSMTOutbound
	CDRTrafficTypeSMSMTInbound
	CDRTrafficTypeSMSMOInbound
	CDRTrafficTypeHTTPMNP
)

var trafficType = map[CDRTrafficType]string{
	CDRTrafficTypeUndefined:     "undefined",
	CDRTrafficTypeSMSMTOutbound: "SMSMT Outbound",
	CDRTrafficTypeSMSMTInbound:  "SMSMT Inbound",
	CDRTrafficTypeSMSMOInbound:  "SMSMO Inbound",
	CDRTrafficTypeHTTPMNP:       "HTTP MNP",
}

type CDRMsuType int

const (
	CDRMsuTypeUndefined CDRMsuType = iota - 1
	_
	CDRMsuTypeItu
	CDRMsuTypeAnsi
)

var msuType = map[CDRMsuType]string{
	CDRMsuTypeUndefined: "undefined",
	CDRMsuTypeItu:       "ITU",
	CDRMsuTypeAnsi:      "ANSI",
}

type CDRRoutingType int

const (
	CDRRoutingTypeUndefined CDRRoutingType = iota - 1
	_
	CDRRoutingTypeFixed
	CDRRoutingTypeMSISDNPrefix
	CDRRoutingTypeOBRPrefix
)

var routingType = map[CDRRoutingType]string{
	CDRRoutingTypeUndefined:    "undefined",
	CDRRoutingTypeFixed:        "Fixed Routing",
	CDRRoutingTypeMSISDNPrefix: "Routing by MSISDN Prefix",
	CDRRoutingTypeOBRPrefix:    "Origin-based Routing by Prefix",
}

type SS7CDR struct {
	MsgID                 string              `json:"msgid"`
	MsgIDTransaction      string              `json:"msgidtran"`
	MsgIDExternal         string              `json:"msgidext"`
	MsgIDProxy            string              `json:"msgidpro"`
	MsgDT                 time.Time           `json:"msgdt"`
	MsgType               CDRMessageType      `json:"msgtype"`
	MsgKind               CDRMessageKind      `json:"msgkind"`
	TrafficDirection      CDRTrafficDirection `json:"trafdir"`
	TrafficType           CDRTrafficType      `json:"traftype"`
	SourceIP              string              `json:"srcip"`
	SourcePort            int                 `json:"srcport"`
	DestinationIP         string              `json:"destip"`
	DestinationPort       int                 `json:"destport"`
	ResultCode            *int                `json:"rescode,omitempty"`
	ResultCodeValid       bool                `json:"rescode_valid"`
	ResultStatus          string              `json:"resstat"`
	MsgDTExpiry           *time.Time          `json:"msgdtexp,omitempty"`
	MsgDTExpiryValid      bool                `json:"msgdtexp_valid"`
	SenderOA              string              `json:"sendoa"`
	DestinationDA         string              `json:"destda"`
	RetryPattern          *int                `json:"retrypat,omitempty"`
	RetryPatternValid     bool                `json:"retrypat_valid"`
	RetryErrorCode        *int                `json:"retryerr,omitempty"`
	RetryErrorCodeValid   bool                `json:"retryerr_valid"`
	RetryAttempt          int                 `json:"retryatt"`
	EsmClass              *int                `json:"esmclass,omitempty"`
	EsmClassValid         bool                `json:"esmclass_valid"`
	Pid64                 *int                `json:"pid64,omitempty"`
	Pid64Valid            bool                `json:"pid64_valid"`
	DataCoding            *int                `json:"datacode,omitempty"`
	DataCodingValid       bool                `json:"datacode_valid"`
	MsgTextLen            *int                `json:"msgtextlen,omitempty"`
	MsgTextLenValid       bool                `json:"msgtextlen_valid"`
	UdhExist              *bool               `json:"udh,omitempty"`
	UdhExistValid         bool                `json:"udh_valid"`
	UdhMsgRefNum          *int                `json:"udhref,omitempty"`
	UdhMsgRefNumValid     bool                `json:"udhref_valid"`
	UdhMsgTotalNum        *int                `json:"udhtot,omitempty"`
	UdhMsgTotalNumValid   bool                `json:"udhtot_valid"`
	UdhMsgPartNum         *int                `json:"udhpart,omitempty"`
	UdhMsgPartNumValid    bool                `json:"udhpart_valid"`
	DlrError              *int                `json:"dlrerr,omitempty"`
	DlrErrorValid         bool                `json:"dlrerr_valid"`
	DlrStatus             string              `json:"dlrstat"`
	TransformationRuleID  int                 `json:"trruleid"`
	SigtranLinkID         int                 `json:"siglinkid"`
	SigtranLinkName       string              `json:"siglinkname"`
	MsuType               CDRMsuType          `json:"msutype"`
	OPC                   int                 `json:"opc"`
	DPC                   int                 `json:"dpc"`
	SccpCarrier           string              `json:"sccpcarrier"`
	SccpClgPaGt           string              `json:"sccpclggt"`
	SccpClgPaTt           *int                `json:"sccpclgtt,omitempty"`
	SccpClgPaTtValid      bool                `json:"sccpclgtt_valid"`
	SccpClgPaSsn          int                 `json:"sccpclgssn"`
	SccpCldPaGt           string              `json:"sccpcldgt"`
	SccpCldPaTt           *int                `json:"sccpcldtt,omitempty"`
	SccpCldPaTtValid      bool                `json:"sccpcldtt_valid"`
	SccpCldPaSsn          int                 `json:"sccpcldssn"`
	TcapID                *int                `json:"tcapid,omitempty"`
	TcapIdValid           bool                `json:"tcapid_valid"`
	MapOpCode             *int                `json:"mapopcode,omitempty"`
	MapOpCodeValid        bool                `json:"mapopcode_valid"`
	MapSCAddr             string              `json:"mapscaddr"`
	MapMscGt              string              `json:"mapmscgt"`
	MapImsi               string              `json:"mapimsi"`
	CustomerAccountID     int                 `json:"custid"`
	CustomerAccountName   string              `json:"custna"`
	SupplierAccountID     int                 `json:"suppid"`
	SupplierAccountName   string              `json:"suppna"`
	DestinationCountryID  int                 `json:"dcouid"`
	DestinationCountry    string              `json:"dcouna"`
	DestinationOperatorID int                 `json:"dopeid"`
	DestinationOperator   string              `json:"dopena"`
	RoutingType           int8                `json:"routty"`
	MsgData               string              `json:"msgdata"`
	MsgDataBin            string              `json:"msgdatabin"`
	UdhData               string              `json:"udhdata"`
	UdhDataBin            string              `json:"udhdatabin"`
}

func ConvertSS7CDRToSs7CdrProc(src *SS7CDR) *dao.Ss7CdrProc {
	if src == nil {
		return nil
	}

	srcIP := validateIPv4(src.SourceIP)
	dstIP := validateIPv4(src.DestinationIP)

	mapScentreAddr := UInt128(src.MapSCAddr)
	mapMscGt := UInt128(src.MapMscGt)

	dst := &dao.Ss7CdrProc{
		InternalMsgID: src.MsgID,
		TranMsgID:     src.MsgIDTransaction,
		ExtMsgID:      src.MsgIDExternal,
		ProxyMsgID:    src.MsgIDProxy,

		MsgDtUs: src.MsgDT,
		MsgDate: src.MsgDT.Truncate(24 * time.Hour),

		MsgType:   int8(src.MsgType),
		MsgKind:   int8(src.MsgKind),
		Direction: int8(src.TrafficDirection),
		Type:      int8(src.TrafficType),
		MsuType:   int8(src.MsuType),

		SrcIP:   srcIP,
		SrcPort: uint16(src.SourcePort),
		DstIP:   dstIP,
		DstPort: uint16(src.DestinationPort),

		ResultCode:   int16(checkIntPointer(src.ResultCode)),
		ResultStatus: src.ResultStatus,

		SenderOA:      src.SenderOA,
		DestinationDA: src.DestinationDA,

		OPC: uint32(src.OPC),
		DPC: uint32(src.DPC),

		SccpCarrier:  src.SccpCarrier,
		SccpClgpaGt:  src.SccpClgPaGt,
		SccpClgpaTt:  int16(checkIntPointer(src.SccpClgPaTt)),
		SccpClgpaSsn: int16(src.SccpClgPaSsn),
		SccpCldpaGt:  src.SccpCldPaGt,
		SccpCldpaTt:  int16(checkIntPointer(src.SccpCldPaTt)),
		SccpCldpaSsn: int16(src.SccpCldPaSsn),

		TcapID: uint32(checkIntPointer(src.TcapID)),

		MapScentreAddr: mapScentreAddr,
		MapMscGt:       mapMscGt,
		MapImsi:        src.MapImsi,
		MapOpco:        int16(checkIntPointer(src.MapOpCode)),

		CustomerAccount:   src.CustomerAccountName,
		CustomerAccountID: uint32(src.CustomerAccountID),
		SupplierAccount:   src.SupplierAccountName,
		SupplierAccountID: uint32(src.SupplierAccountID),

		SignallingConnLink:   src.SigtranLinkName,
		SignallingConnLinkID: uint32(src.SigtranLinkID),

		DestinationCountry:    src.DestinationCountry,
		DestinationCountryID:  uint32(src.DestinationCountryID),
		DestinationOperator:   src.DestinationOperator,
		DestinationOperatorID: uint32(src.DestinationOperatorID),

		EsmClass:    int16(checkIntPointer(src.EsmClass)),
		DataCoding:  int16(checkIntPointer(src.DataCoding)),
		Pid64:       int8(checkIntPointer(src.Pid64)),
		MsgTextLen:  int16(checkIntPointer(src.MsgTextLen)),
		Udh:         0,
		MsgRefNum:   uint16(checkIntPointer(src.UdhMsgRefNum)),
		MsgTotalNum: int8(checkIntPointer(src.UdhMsgTotalNum)),
		MsgPartNum:  int8(checkIntPointer(src.UdhMsgPartNum)),

		DlrErr:  int16(checkIntPointer(src.DlrError)),
		DlrStat: src.DlrStatus,

		RetryPattern: int8(checkIntPointer(src.RetryPattern)),
		RetryError:   int16(checkIntPointer(src.RetryErrorCode)),
		RetryAttempt: int16(src.RetryAttempt),

		RoutingType:          int8(src.RoutingType),
		TransformationRuleID: uint16(src.TransformationRuleID),

		MsgData:    src.MsgData,
		MsgDataBin: src.MsgDataBin,
		UdhData:    src.UdhData,
		UdhDataBin: src.UdhDataBin,
	}

	if src.MsgDTExpiry != nil {
		dst.MsgExpiryDt = *src.MsgDTExpiry
	}

	if src.UdhExist != nil {
		if *src.UdhExist {
			dst.Udh = 1
		} else {
			dst.Udh = 0
		}
	} else {
		dst.Udh = -1
	}

	return dst
}

func checkIntPointer(value *int) int {
	if value != nil {
		return *value
	}
	return -1
}

func validateIPv4(ipStr string) string {

	if ipStr == "" {
		return "0.0.0.0"
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "0.0.0.0"
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return "0.0.0.0"
	}

	return ipv4.String()
}

func UInt128(v string) *big.Int {
	if v == "" {
		return big.NewInt(0)
	}

	bigInt := new(big.Int)
	bigInt, ok := bigInt.SetString(v, 10)
	if !ok {
		return big.NewInt(0)
	}

	if bigInt.Sign() < 0 {
		return big.NewInt(0)
	}

	maxUInt128 := new(big.Int)
	maxUInt128.Exp(big.NewInt(2), big.NewInt(128), nil)
	maxUInt128.Sub(maxUInt128, big.NewInt(1))

	if bigInt.Cmp(maxUInt128) > 0 {
		return big.NewInt(0)
	}

	return bigInt
}
