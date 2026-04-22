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

	if err := conn.Exec(context.Background(), getRepoInsQuery(dao.Ss7CdrProc{}), transaction); err != nil {
		return err
	}

	return nil
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
