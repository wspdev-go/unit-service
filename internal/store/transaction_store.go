package store

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unit-service/internal/config"
	"unit-service/logger"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type TransactionStore interface {
	Open() error
	Close() error
	Ping() error
	Conn() clickhouse.Conn
}

type transaction struct {
	cfg  *config.TransactionConfig
	conn clickhouse.Conn
}

func NewTransaction(cfg *config.TransactionConfig) TransactionStore {
	return &transaction{
		cfg: cfg,
	}
}

func (t *transaction) Open() error {

	if t.conn != nil {
		return nil
	}

	if t.cfg == nil {
		return errors.New("config is nil")
	}

	if t.cfg.Host == "" {
		return errors.New("host is empty")
	}

	if t.cfg.Port == 0 {
		return errors.New("port is empty")
	}

	if t.cfg.Database == "" {
		return errors.New("database is empty")
	}

	if t.cfg.Username == "" {
		return errors.New("username is empty")
	}

	if t.cfg.Password == "" {
		return errors.New("password is empty")
	}

	connAdr := fmt.Sprintf("%s:%d", t.cfg.Host, t.cfg.Port)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{connAdr},
		Auth: clickhouse.Auth{
			Database: t.cfg.Database,
			Username: t.cfg.Username,
			Password: t.cfg.Password,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:  time.Duration(t.cfg.DialTimeout) * time.Second,
		MaxOpenConns: t.cfg.MaxOpenConns,
		MaxIdleConns: t.cfg.MaxIdleConns,
		//ConnOpenStrategy: clickhouse.ConnOpenRoundRobin,//Need to use this strategy if we have multiple ClickHouse nodes
	})

	if err != nil {
		logger.Error("open clickhouse connection error: %v", err)
		return err
	}

	t.conn = conn

	// Check if the connection is alive
	if err = t.Ping(); err != nil {
		logger.Error("clickhouse ping error: %v", err)
		_ = t.Close()
		return err
	}

	return nil
}

func (t *transaction) Ping() error {
	if t.conn == nil {
		return fmt.Errorf("transaction store is not open")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.conn.Ping(ctx)

}

func (t *transaction) Close() error {
	if t.conn == nil {
		return nil
	}

	if err := t.conn.Close(); err != nil {
		logger.Error("Failed to close Transaction connection: %s", err.Error())
		return err
	}

	t.conn = nil

	logger.Info("close transaction successfully")

	return nil
}

func (t *transaction) Conn() clickhouse.Conn {
	return t.conn
}
