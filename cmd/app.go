package cmd

import (
	"context"
	"log/slog"
	"marketflow/internal/adapters/cache"
	"marketflow/internal/adapters/postgres"
	"marketflow/internal/adapters/web"
	"marketflow/internal/config"
	"marketflow/internal/infra"
	"os"
	"os/signal"
	"syscall"
)

func RunApp() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	infra.SetUpLogger()

	config, err := config.ReadConfig() // Read config
	if err != nil {
		slog.Error("Config Error", "error", err)
	}

	rdb, err := cache.NewRedisClient(&config) // Connect to Redis
	if err != nil {
		slog.Error("Redis Error", "error", err)
	}

	db := postgres.ConnectDB(&config)
	repository := postgres.NewRepository(db)

	// exchange.RunTCPClients(&config, rdb, false)

	

	err = web.StartServer(&config)
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
