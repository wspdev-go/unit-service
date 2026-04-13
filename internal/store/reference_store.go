package store

import (
	"fmt"
	"unit-service/internal/config"
	"unit-service/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ReferenceStore interface {
	Open() bool
	Close() error
	Ping() error
}

type reference struct {
	cfg  *config.ReferenceConfig
	dsn  string
	Conn *gorm.DB
}

func NewReference(cfg *config.ReferenceConfig) ReferenceStore {

	if cfg == nil {
		return &reference{}
	}

	dbPort := fmt.Sprintf("%d", cfg.Port)

	dsn := cfg.Username + ":" + cfg.Password +
		"@tcp" + "(" + cfg.Host + ":" + dbPort + ")/" + cfg.Database +
		"?" + "parseTime=true&loc=Local"

	return &reference{
		cfg: cfg,
		dsn: dsn,
	}
}

func (r *reference) Open() bool {
	if r.Conn != nil {
		return true
	}

	if r.cfg == nil {
		logger.Error("reference config is nil")
		return false
	}

	if r.cfg.Host == "" {
		logger.Error("reference host is empty")
		return false
	}

	if r.cfg.Port == 0 {
		logger.Error("reference port is empty")
		return false
	}

	if r.cfg.Database == "" {
		logger.Error("reference database is empty")
		return false
	}

	if r.cfg.Username == "" {
		logger.Error("reference username is empty")
		return false
	}

	if r.cfg.Password == "" {
		logger.Error("reference password is empty")
		return false
	}

	conn, err := gorm.Open(mysql.Open(r.dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Can't open gorm: %s", err)
		return false
	}

	r.Conn = conn

	if err = r.Ping(); err != nil {
		logger.Error("can't ping reference store: %s", err)
		_ = r.Close()
		return false
	}

	return true
}

func (r *reference) Ping() error {
	if r.Conn == nil {
		return fmt.Errorf("reference store is not open")
	}

	sqlDB, err := r.Conn.DB()
	if err != nil {
		logger.Error("Can't get sqlDB: %s", err)
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		logger.Error("can't ping reference store: %s", err)
		_ = r.Close()
		return err
	}

	return nil

}

func (r *reference) Close() error {
	if r.Conn == nil {
		return nil
	}

	dbConn, err := r.Conn.DB()
	if err != nil {
		logger.Error("Can't close reference store: %s", err)
		return err
	}

	if err = dbConn.Close(); err != nil {
		logger.Error("can't close reference store: %s", err)
		return err
	}
	r.Conn = nil

	logger.Info("Reference store closed successfully")
	return nil
}
