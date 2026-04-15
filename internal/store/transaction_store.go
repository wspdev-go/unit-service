package store

import (
	"context"
	"fmt"
	"time"
	"unit-service/internal/config"
	"unit-service/logger"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type TransactionStore interface {
	Open() bool
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

func (t *transaction) Open() bool {

	if t.conn != nil {
		return true
	}

	if t.cfg == nil {
		logger.Error("transaction config is nil")
		return false
	}

	if t.cfg.Host == "" {
		logger.Error("host is empty")
		return false
	}

	if t.cfg.Port == 0 {
		logger.Error("port is empty")
		return false
	}

	if t.cfg.Database == "" {
		logger.Error("dbName is empty")
		return false
	}

	if t.cfg.Username == "" {
		logger.Error("username is empty")
		return false
	}

	if t.cfg.Password == "" {
		logger.Error("password is empty")
		return false
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
		return false
	}

	t.conn = conn

	// Check if the connection is alive
	if err = t.Ping(); err != nil {
		logger.Error("clickhouse ping error: %v", err)
		_ = t.Close()
		return false
	}

	return true
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
