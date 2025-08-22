package cmd

import (
	"log/slog"
	"marketflow/internal/adapters/cache"
	"marketflow/internal/adapters/exchange"
	"marketflow/internal/adapters/web"
	"marketflow/internal/config"
	"marketflow/internal/infra"
)

func RunApp() {
	infra.SetUpLogger()

	config, err := config.ReadConfig() // Read config
	if err != nil {
		slog.Error("Config Error", "error", err)
	}

	rdb, err := cache.NewRedisClient(&config) // Connect to Redis
	if err != nil {
		slog.Error("Redis Error", "error", err)
	}

	// exchange.RunTCPClients(&config, rdb, false)

	err = web.StartServer(&config)
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
