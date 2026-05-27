/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unit-service/internal/config"
	"unit-service/internal/store"
	"unit-service/internal/usecase"
	"unit-service/logger"
)

type App struct {
	Config *config.Config
	Store  store.Store
}

func NewApp(configPath string) (*App, error) {
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger.Info("Application initialized with configuration: %s", configPath)

	st := store.NewStore(cfg)

	return &App{
		Config: cfg,
		Store:  st,
	}, nil
}

func (a *App) RunQueuePipeline() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	requiredDeps := []Dependency{
		DepReference,
		DepQueue,
		DepTransaction,
	}

	if err := a.RunDependency(requiredDeps...); err != nil {
		return err
	}

	defer func() {
		if err := a.StopDependency(requiredDeps...); err != nil {
			logger.Error("close connections: %v", err)
		}
	}()

	uc := usecase.NewUsecase(a.Store)
	return uc.RunQueuePipeline(ctx)
}
