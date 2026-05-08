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
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unit-service/internal/config"
	"unit-service/internal/repository"
	"unit-service/internal/store"
	"unit-service/internal/usecase"
	"unit-service/logger"

	"golang.org/x/sync/errgroup"
)

type App struct {
	Config *config.Config
	Store  store.Store
	Repo   repository.Repository
}

func NewApp(configPath string) (*App, error) {
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger.Info("Application initialized with configuration: %s", configPath)

	st := store.NewStore(cfg)

	repo := repository.NewRepository(st)

	return &App{
		Config: cfg,
		Store:  st,
		Repo:   repo,
	}, nil
}

func (a *App) RunApp() {
	// Start the application logic here
	// This function can be used to run the main application loop, handle requests, etc.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	uc := usecase.NewUsecase(a.Repo)

	refUc, err := uc.GetReferenceUsecase()
	if err != nil {
		logger.Error("error initializing reference use case: %v", err)
		return
	}

	if err = refUc.Run(ctx); err != nil {
		logger.Error("error loading reference data: %v", err)
		return
	}

	transactionUc, err := uc.GetTransactionUsecase()
	if err != nil {
		logger.Error("error initializing transaction use case: %v", err)
		return
	}

	queueUc, err := uc.GetQueueUsecase()
	if err != nil {
		logger.Error("Error initializing queue use case: %v", err)
		return
	}

	g, groupCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return transactionUc.Run(groupCtx)
	})

	g.Go(func() error {
		return queueUc.Run(groupCtx)
	})

	if err = g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger.Info("application stopped")
			return
		}

		logger.Error("error running application workers: %v", err)
		return
	}
}
