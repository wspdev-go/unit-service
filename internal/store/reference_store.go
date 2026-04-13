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
}

type reference struct {
	cfg  *config.ReferenceConfig
	dsn  string
	Conn *gorm.DB
}

func NewReference(cfg *config.ReferenceConfig) ReferenceStore {

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
	var err error

	if r.Conn != nil {
		return true
	}

	r.Conn, err = gorm.Open(mysql.Open(r.dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Can't open gorm: %s", err)
		return false
	}
	return true
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
	err = dbConn.Close()
	if err != nil {
		logger.Error("Can't close reference store: %s", err)
	}
	logger.Info("Reference store closed successfully")
	return nil
}
