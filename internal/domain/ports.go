package domain

import (
	"context"
	"marketflow/internal/config"
	"net/http"
	"time"
)

type ExchangeClient interface {
	StartLiveMode(ctx context.Context) <-chan PriceTick
	StartTestMode(ctx context.Context) <-chan PriceTick
	Stop()
}

type RedisClient interface {
	StoreTick(ctx context.Context, exchange, pair string, price float64) error
	ProcessLastMinute(ctx context.Context) []*MinuteAgg

	// --- latest ---
	GetLatestPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error)

	// --- highest ---
	GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)

	// --- lowest ---
	GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)

	// --- average ---
	GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)
}

type Repository interface {
	// Infra
	Ping(ctx context.Context) error

	// Write-side (minute aggregates)
	StoreMinAgg(ctx context.Context, aggs []*MinuteAgg) error

	// -------- HIGHEST --------
	GetHighestPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error)
	GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)

	// -------- LOWEST ---------
	GetLowestPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetLowestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error)
	GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)

	// -------- AVERAGE --------
	GetAvgPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetAvgPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error)
	GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)
}

type WindowStore interface {
	// Запись и снапшот
	Add(msg PriceTick)
	Snapshot() map[Key][]PriceTick

	// Latest
	GetLatestPriceBySymbol(symbol string) (float64, error)
	GetLatestPriceBySymbolAndExchange(symbol, exchange string) (float64, error)

	// Highest
	GetHighestPriceBySymbolAndPeriod(symbol string, period time.Duration) (float64, error)
	GetHighestPriceBySymbolAndPeriodAndExchange(symbol, exchange string, period time.Duration) (float64, error)

	// Lowest
	GetLowestPriceBySymbol(symbol string) (float64, error)
	GetLowestPriceBySymbolAndExchange(symbol, exchange string) (float64, error)
	GetLowestPriceBySymbolAndPeriod(symbol string, period time.Duration) (float64, error)
	GetLowestPriceBySymbolAndPeriodAndExchange(symbol, exchange string, period time.Duration) (float64, error)

	// Average
	GetAvgPriceBySymbol(symbol string) (float64, error)
	GetAvgPriceBySymbolAndExchange(symbol, exchange string) (float64, error)
	GetAvgPriceBySymbolAndPeriod(symbol string, period time.Duration) (float64, error)
	GetAvgPriceBySymbolAndPeriodAndExchange(symbol, exchange string, period time.Duration) (float64, error)
}

type DataProcessingService interface {
	StartWorkers(ctx context.Context, in <-chan PriceTick)
	StopWorkers()
	ExchangesHealth(within time.Duration, names []string) map[string]bool
}

type ServerHandler interface {
	StartServer(ctx context.Context, config *config.Config) error
	RegisterRouter(mux *http.ServeMux)
}

type MarketDataService interface {
	GetLatestPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error)

	GetHighestPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error)
	GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)

	GetLowestPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetLowestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error)
	GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)

	GetAvgPriceBySymbol(ctx context.Context, symbol string) (float64, error)
	GetAvgPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error)
	GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error)
	GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) (float64, error)
}

type SystemService interface {
	SwitchToTestMode() error
	SwitchToLiveMode() error
	GetCurrentMode(ctx context.Context) (string, error)
	IsLiveMode() bool
	Shutdown(ctx context.Context) error
	// Health check
	GetSystemHealth(ctx context.Context) (SystemHealth, error)

	// Dependency status
	CheckRedisHealth(ctx context.Context) bool
	CheckPostgresHealth(ctx context.Context) bool
	CheckExchangeHealth(ctx context.Context) map[string]bool
}
