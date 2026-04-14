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

func (a *App) OpenConnections() error {
	transaction := a.Store.GetTransaction()
	if transaction == nil {
		return fmt.Errorf("transaction store is nil")
	}
	if !transaction.Open() {
		return fmt.Errorf("open transaction store")
	}
	logger.Info("Transaction store is open and ready to use.")

	queue := a.Store.GetQueue()
	if queue == nil {
		return fmt.Errorf("queue store is nil")
	}
	if !queue.Open() {
		return fmt.Errorf("open queue store")
	}
	logger.Info("Queue store is open and ready to use.")

	reference := a.Store.GetReference()
	if reference == nil {
		return fmt.Errorf("reference store is nil")
	}
	if !reference.Open() {
		return fmt.Errorf("open reference store")
	}
	logger.Info("Reference store is open and ready to use.")

	return nil
}

func (a *App) CloseConnections() error {
	var firstErr error

	if transaction := a.Store.GetTransaction(); transaction != nil {
		if err := transaction.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close transaction store: %w", err)
		}
	}

	if queue := a.Store.GetQueue(); queue != nil {
		if err := queue.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close queue store: %w", err)
		}
	}

	if reference := a.Store.GetReference(); reference != nil {
		if err := reference.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close reference store: %w", err)
		}
	}

	return firstErr
}

func RunApp() {
	// Start the application logic here
	// This function can be used to run the main application loop, handle requests, etc.
}
