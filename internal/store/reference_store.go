/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package store

import (
	"errors"
	"fmt"
	"unit-service/internal/config"
	"unit-service/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ReferenceStore interface {
	Open() error
	Close() error
	Ping() error
	DB() *gorm.DB
}

type reference struct {
	cfg *config.ReferenceConfig
	dsn string
	db  *gorm.DB
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

func (r *reference) Open() error {
	if r.db != nil {
		return nil
	}

	if r.cfg == nil {
		logger.Error("reference config is nil")
		return errors.New("reference config is nil")
	}

	if r.cfg.Host == "" {
		logger.Error("reference host is empty")
		return errors.New("reference host is empty")
	}

	if r.cfg.Port == 0 {
		logger.Error("reference port is empty")
		return errors.New("reference port is empty")
	}

	if r.cfg.Database == "" {
		logger.Error("reference database is empty")
		return errors.New("reference database is empty")
	}

	if r.cfg.Username == "" {
		logger.Error("reference username is empty")
		return errors.New("reference username is empty")
	}

	if r.cfg.Password == "" {
		logger.Error("reference password is empty")
		return errors.New("reference password is empty")
	}

	gormDB, err := gorm.Open(mysql.Open(r.dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Can't open gorm: %s", err)
		return err
	}

	r.db = gormDB

	if err = r.Ping(); err != nil {
		logger.Error("can't ping reference store: %s", err)
		_ = r.Close()
		return err
	}

	return nil
}

func (r *reference) Ping() error {
	if r.db == nil {
		return fmt.Errorf("reference store is not open")
	}

	sqlDB, err := r.db.DB()
	if err != nil {
		logger.Error("Can't get sqlDB: %s", err)
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		logger.Error("can't ping reference store: %s", err)
		return err
	}

	return nil

}

func (r *reference) Close() error {
	if r.db == nil {
		return nil
	}

	dbConn, err := r.db.DB()
	if err != nil {
		logger.Error("Can't close reference store: %s", err)
		return err
	}

	if err = dbConn.Close(); err != nil {
		logger.Error("can't close reference store: %s", err)
		return err
	}
	r.db = nil

	logger.Info("Reference store closed successfully")
	return nil
}

func (r *reference) DB() *gorm.DB {
	return r.db
}
