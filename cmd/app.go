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
	"time"
)

func RunApp() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	infra.SetUpLogger()

	config, err := config.ReadConfig()
	if err != nil {
		slog.Error("Config Error", "error", err)
		os.Exit(1)
	}

	rdb, err := cache.NewRedisClient(&config)
	if err != nil {
		slog.Error("Redis Error", "error", err)
	}

	db := postgres.ConnectDB(&config)
	repository := postgres.NewRepository(db)
	exchangeClient := exchange.NewExchangeClient(&config, 5*time.Second)
	
	dataProcessingService := app.NewDataProcessingService(rdb, repository)
	mds := app.NewMarketDataService(rdb, repository)
	ss := app.NewModeService(exchangeClient, dataProcessingService, rdb, repository, &config)
	
	rootCtx := context.Background()

	if config.Mode == "live" {
		if err := ss.SwitchToLiveMode(rootCtx); err != nil {
			slog.Error("Failed to switch to live mode", "error", err)
			return
		}
	} else {
		if err := ss.SwitchToTestMode(rootCtx); err != nil {
			slog.Error("Failed to switch to test mode", "error", err)
			return
		}
	}
	
	handler := web.NewHandler(mds, ss)
	err = handler.StartServer(ctx, &config)
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}