package app

import (
	"unit-service/internal/config"
	"unit-service/internal/store"
	"unit-service/logger"
)

func InitApp(configPath string) {
	// Initialize the application with the provided configuration path
	// This function can be used to set up necessary resources, configurations, and dependencies
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		panic(err)
	}

	logger.Info("Application initialized with configuration: %s", configPath)

	store := store.NewStore(cfg)

	storeTransaction := store.GetTransaction()
	if storeTransaction != nil {
		panic("Transaction store should be nil at initialization")
	}

	storeQueue := store.GetQueue()
	if storeQueue != nil {
		panic("Queue store should be nil at initialization")
	}

	storeReference := store.GetReference()
	if storeReference == nil {
		panic("Reference store should not be nil at initialization")
	}

}

func RunApp() {
	// Start the application logic here
	// This function can be used to run the main application loop, handle requests, etc.
}
