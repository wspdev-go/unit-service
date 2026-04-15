package app

import (
	"fmt"
	"unit-service/internal/config"
	"unit-service/internal/repository"
	"unit-service/internal/store"
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

func RunApp() {
	// Start the application logic here
	// This function can be used to run the main application loop, handle requests, etc.
}
