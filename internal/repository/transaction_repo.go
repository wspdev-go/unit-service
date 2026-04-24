package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unit-service/internal/model/dao"
	"unit-service/internal/store"
)

type TransactionRepo interface {
	PushBatch() error
	PutTransaction(transaction *dao.Ss7CdrProc) error
}

type transactionRepo struct {
	store store.TransactionStore
}

func NewTransactionRepo(store store.TransactionStore) TransactionRepo {
	return &transactionRepo{
		store: store,
	}
}

func (repo *transactionRepo) PutTransaction(transaction *dao.Ss7CdrProc) error {
	conn := repo.store.Conn()
	if conn == nil {
		return errors.New("conn is nil")
	}

	query := getRepoInsQuery(dao.Ss7CdrProc{})

	if err := conn.Exec(context.Background(), query, getCdrFields(transaction)...); err != nil {
		return err
	}

	return nil
}

func getCdrFields(cdr *dao.Ss7CdrProc) []any {
	return []any{
		cdr.MsgDate,
		cdr.MsgDtUs,
		cdr.MsgExpiryDt,

		cdr.ExtMsgID,
		cdr.ProxyMsgID,
		cdr.InternalMsgID,
		cdr.TranMsgID,

		// IP and Port information
		cdr.SrcIP,
		cdr.SrcPort,
		cdr.DstIP,
		cdr.DstPort,

		// Message types and kinds
		cdr.MsgType,
		cdr.MsgKind,
		cdr.MsuType,
		cdr.Type,

		// Direction
		cdr.Direction,

		// Result information
		cdr.ResultCode,
		cdr.ResultStatus,

		// Message addresses
		cdr.SenderOA,
		cdr.DestinationDA,

		cdr.OPC,
		cdr.DPC,

		cdr.SccpCarrier,
		cdr.SccpClgpaGt,
		cdr.SccpClgpaTt,
		cdr.SccpClgpaSsn,
		cdr.SccpCldpaGt,
		cdr.SccpCldpaTt,
		cdr.SccpCldpaSsn,

		cdr.TcapID,

		cdr.MapScentreAddr,
		cdr.MapMscGt,
		cdr.MapImsi,
		cdr.MapOpco,

		cdr.CustomerAccount,
		cdr.CustomerAccountID,
		cdr.SupplierAccount,
		cdr.SupplierAccountID,

		cdr.SignallingConnLink,
		cdr.SignallingConnLinkID,

		cdr.DestinationCountry,
		cdr.DestinationCountryID,
		cdr.DestinationOperator,
		cdr.DestinationOperatorID,

		cdr.EsmClass,
		cdr.DataCoding,
		cdr.Pid64,
		cdr.MsgTextLen,
		cdr.Udh,
		cdr.MsgRefNum,
		cdr.MsgTotalNum,
		cdr.MsgPartNum,

		// DLR information
		cdr.DlrErr,
		cdr.DlrStat,

		// Retry information
		cdr.RetryPattern,
		cdr.RetryError,
		cdr.RetryAttempt,

		cdr.RoutingType,
		cdr.TransformationRuleID,

		cdr.MsgData,
		cdr.MsgDataBin,
		cdr.UdhData,
		cdr.UdhDataBin,
	}
}

func getRepoInsQuery(obj dao.Ss7CdrProc) string {
	t := reflect.TypeOf(obj)

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

func (repo *transactionRepo) PushBatch() error {
	return nil
}
