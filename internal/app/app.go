package app

import (
	"fmt"
	"unit-service/internal/config"
	"unit-service/internal/repository"
	"unit-service/internal/store"
	"unit-service/logger"
)

type Dependency string

const (
	DepReference   Dependency = "reference"
	DepQueue       Dependency = "queue"
	DepTransaction Dependency = "transaction"
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

func (a *App) OpenConnections(required ...Dependency) error {
	for _, dep := range required {
		switch dep {
		case DepTransaction:
			if err := a.openTransaction(); err != nil {
				return err
			}
		case DepQueue:
			if err := a.openQueue(); err != nil {
				return err
			}
		case DepReference:
			if err := a.openReference(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown dependency: %s", dep)
		}
	}

	return nil
}

func (a *App) openTransaction() error {
	transaction := a.Store.GetTransaction()
	if transaction == nil {
		return fmt.Errorf("transaction store is nil")
	}
	if err := transaction.Open(); err != nil {
		return fmt.Errorf("open transaction store")
	}

	logger.Info("Transaction store is open and ready to use.")
	return nil
}

func (a *App) openQueue() error {
	queue := a.Store.GetQueue()
	if queue == nil {
		return fmt.Errorf("queue store is nil")
	}
	if err := queue.Open(); err != nil {
		return err
	}

	logger.Info("Queue store is open and ready to use.")
	return nil
}

func (a *App) openReference() error {
	reference := a.Store.GetReference()
	if reference == nil {
		return fmt.Errorf("reference store is nil")
	}
	if err := reference.Open(); err != nil {
		return err
	}

	logger.Info("Reference store is open and ready to use.")
	return nil
}

func (a *App) CloseConnections(required ...Dependency) error {
	var firstErr error

	for _, dep := range required {
		switch dep {
		case DepTransaction:
			if transaction := a.Store.GetTransaction(); transaction != nil {
				if err := transaction.Close(); err != nil && firstErr == nil {
					firstErr = fmt.Errorf("close transaction store: %w", err)
				}
			}
		case DepQueue:
			if queue := a.Store.GetQueue(); queue != nil {
				if err := queue.Close(); err != nil && firstErr == nil {
					firstErr = fmt.Errorf("close queue store: %w", err)
				}
			}
		case DepReference:
			if reference := a.Store.GetReference(); reference != nil {
				if err := reference.Close(); err != nil && firstErr == nil {
					firstErr = fmt.Errorf("close reference store: %w", err)
				}
			}
		default:
			if firstErr == nil {
				firstErr = fmt.Errorf("unknown dependency: %s", dep)
			}
		}
	}

	return firstErr
}

func RunApp() {
	// Start the application logic here
	// This function can be used to run the main application loop, handle requests, etc.
}
