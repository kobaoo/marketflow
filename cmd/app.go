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
	config.ParseFlags()
	infra.SetUpLogger()

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ReadConfig()
	if err != nil {
		slog.Error("Config Error", "error", err)
		os.Exit(1)
	}

	rdb, err := cache.NewRedisClient(&cfg)
	if err != nil {
		slog.Error("Redis Error", "error", err)
	}

	db := postgres.ConnectDB(&cfg)
	repository := postgres.NewRepository(db)

	exchangeClient := exchange.NewExchangeClient(&cfg, 5*time.Second)

	window := app.NewWindowStore()
	dataProcessingService := app.NewDataProcessingService(rdb, repository, window)
	marketDataService := app.NewMarketDataService(rdb, repository, window)
	sysService := app.NewModeService(rootCtx, exchangeClient, dataProcessingService, rdb, repository, &cfg)

	switch cfg.Mode {
	case "live":
		if err := sysService.SwitchToLiveMode(); err != nil {
			slog.Error("Failed to switch to live mode", "error", err)
			return
		}
	default:
		if err := sysService.SwitchToTestMode(); err != nil {
			slog.Error("Failed to switch to test mode", "error", err)
			return
		}
	}

	handler := web.NewHandler(marketDataService, sysService)
	go func() {
		if err := handler.StartServer(rootCtx, &cfg); err != nil {
			slog.Error("Error starting server", "error", err)
		}
	}()

	<-rootCtx.Done()
	slog.Info("Interrupt received, shutting down...")

	_ = sysService.Shutdown(context.Background())

	slog.Info("Shutdown complete")
}
