package cmd

import (
	"context"
	"log/slog"
	"marketflow/internal/adapters/cache"
	"marketflow/internal/adapters/exchange"
	"marketflow/internal/adapters/postgres"
	"marketflow/internal/adapters/web"
	"marketflow/internal/app"
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

	exchangeClient := exchange.NewExchangeClient(&config)
	if config.Mode == "live" {
		messages := exchangeClient.StartLiveMode(ctx)
		dataProcessingService := app.NewDataProcessingService(rdb, repository)
		dataProcessingService.StartWorkers(ctx, messages)
	} else {
		messages := exchangeClient.StartTestMode(ctx)
		dataProcessingService := app.NewDataProcessingService(rdb, repository)
		dataProcessingService.StartWorkers(ctx, messages)
	}

	mds := app.NewMarketDataService(rdb, repository)
	handler := web.NewHandler(mds)
	err = web.StartServer(ctx, &config, handler)
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
