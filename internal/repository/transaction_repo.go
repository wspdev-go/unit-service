package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unit-service/internal/model/dao"
	"unit-service/internal/store"

	"github.com/ClickHouse/clickhouse-go/v2"
)

const batchSize = 1000

type TransactionRepo interface {
	PushBatchTransaction(ctx context.Context) error
	PutBatch(transaction *dao.Ss7CdrProc) error
	PutTransaction(transaction *dao.Ss7CdrProc) error
	FlushCh() <-chan struct{}
	GetConnValid() bool
	SetConnValid(valid bool)
	ConnRecovery(ctx context.Context) error
}

type transactionRepo struct {
	conn        clickhouse.Conn
	store       store.TransactionStore
	batchBuff   []dao.Ss7CdrProc
	mu          sync.Mutex
	flushCh     chan struct{}
	isConnValid atomic.Bool
}

func NewTransactionRepo(store store.TransactionStore) (TransactionRepo, error) {
	conn := store.Conn()

	if conn == nil {
		return nil, errors.New("clickhouse connection is nil")
	}

	return &transactionRepo{
		conn:      conn,
		store:     store,
		batchBuff: make([]dao.Ss7CdrProc, 0, batchSize), // Initialize batch buffer with a reasonable capacity
		flushCh:   make(chan struct{}, 1),
	}, nil
}

func (repo *transactionRepo) PutTransaction(transaction *dao.Ss7CdrProc) error {
	if repo.conn == nil {
		return errors.New("conn is nil")
	}

	query := getRepoInsQuery(dao.Ss7CdrProc{})

	if err := repo.conn.Exec(context.Background(), query, getCdrFields(transaction)...); err != nil {
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

func (repo *transactionRepo) PutBatch(transaction *dao.Ss7CdrProc) error {
	if transaction == nil {
		return errors.New("transaction is nil")
	}

	repo.mu.Lock()
	repo.batchBuff = append(repo.batchBuff, *transaction)
	if len(repo.batchBuff) >= batchSize {
		select {
		case repo.flushCh <- struct{}{}:
		default:
		}
	}
	repo.mu.Unlock()

	return nil
}

func (repo *transactionRepo) FlushCh() <-chan struct{} {
	return repo.flushCh
}

func (repo *transactionRepo) PushBatchTransaction(ctx context.Context) error {
	repo.mu.Lock()

	if len(repo.batchBuff) == 0 {
		repo.mu.Unlock()
		return nil
	}

	if !repo.GetConnValid() {
		repo.mu.Unlock()
		return nil
	}

	buff := make([]dao.Ss7CdrProc, len(repo.batchBuff))
	copy(buff, repo.batchBuff)
	repo.batchBuff = make([]dao.Ss7CdrProc, 0, 3*batchSize)

	repo.mu.Unlock()

	return repo.runPush(ctx, buff)
}

func (repo *transactionRepo) runPush(ctx context.Context, buff []dao.Ss7CdrProc) error {
	batch, err := repo.conn.PrepareBatch(ctx, getRepoInsQuery(dao.Ss7CdrProc{}))
	if err != nil {
		repo.restoreFailedBatch(buff)
		return err
	}

	for _, cdr := range buff {
		if err = batch.Append(getCdrFields(&cdr)...); err != nil {
			_ = batch.Abort()
			repo.restoreFailedBatch(buff)
			return err
		}
	}

	if err = batch.Send(); err != nil {
		_ = batch.Abort()
		repo.restoreFailedBatch(buff)
		return err
	}

	batch = nil

	return nil
}

func (repo *transactionRepo) restoreFailedBatch(buff []dao.Ss7CdrProc) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.batchBuff = append(buff, repo.batchBuff...)
}

func (repo *transactionRepo) GetConnValid() bool {
	return repo.isConnValid.Load()
}

func (repo *transactionRepo) SetConnValid(valid bool) {
	repo.isConnValid.Store(valid)
}

func (repo *transactionRepo) ConnRecovery(ctx context.Context) error {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if repo.GetConnValid() {
				return nil
			}

			if err := repo.store.Ping(); err == nil {
				repo.SetConnValid(true)
				return nil
			}
		}
	}
}
