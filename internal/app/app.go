package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unit-service/internal/config"
	"unit-service/internal/repository"
	"unit-service/internal/store"
	"unit-service/internal/usecase"
	"unit-service/logger"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	uc := usecase.NewUsecase(a.Repo)

	queueUc, err := uc.GetQueueUsecase()
	if err != nil {
		logger.Error("Error initializing queue use case: %v", err)
		return
	}

	if err = queueUc.Run(ctx); err != nil {
		logger.Error("Error running queue use case: %v", err)
		return
	}
}
