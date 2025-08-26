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
    // Логгер
    infra.SetUpLogger()

    // Root context с перехватом сигналов
    rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    cfg, err := config.ReadConfig()
    if err != nil {
        slog.Error("Config Error", "error", err)
        os.Exit(1)
    }

    // Redis
    rdb, err := cache.NewRedisClient(&cfg)
    if err != nil {
        slog.Error("Redis Error", "error", err)
    }

    // Postgres
    db := postgres.ConnectDB(&cfg)
    repository := postgres.NewRepository(db)

    // Источники (биржи)
    exchangeClient := exchange.NewExchangeClient(&cfg, 5*time.Second)

    // Сервисы
    dataProcessingService := app.NewDataProcessingService(rdb, repository)
    marketDataService := app.NewMarketDataService(rdb, repository)
    sysService := app.NewModeService(exchangeClient, dataProcessingService, rdb, repository, &cfg)

    // Стартовый режим
    switch cfg.Mode {
    case "live":
        if err := sysService.SwitchToLiveMode(rootCtx); err != nil {
            slog.Error("Failed to switch to live mode", "error", err)
            return
        }
    default:
        if err := sysService.SwitchToTestMode(rootCtx); err != nil {
            slog.Error("Failed to switch to test mode", "error", err)
            return
        }
    }

    // HTTP-сервер (должен уважать завершение по ctx.Done() внутри)
    handler := web.NewHandler(marketDataService, sysService)
    go func() {
        if err := handler.StartServer(rootCtx, &cfg); err != nil {
            slog.Error("Error starting server", "error", err)
        }
    }()

    // Ждём сигнал
    <-rootCtx.Done()
    slog.Info("Interrupt received, shutting down...")

    // Корректно гасим сервисы
    _ = sysService.Shutdown(context.Background())

    slog.Info("Shutdown complete")
}
