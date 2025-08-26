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
	StoreTick(ctx context.Context, exchange, pair string, price float64)
	ProcessLastMinute(ctx context.Context) []*MinuteAgg

	GetLatestPriceByPattern(ctx context.Context, pattern string) float64

	GetLatestPriceBySymbol(ctx context.Context, symbol string) float64
	GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64

	GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64

	GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64

	GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64
}

type Repository interface {
	StoreMinAgg(ctx context.Context, aggs []*MinuteAgg)

	GetHighestPriceBySymbol(ctx context.Context, symbol string) float64
	GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64
	GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64

	GetLowestPriceBySymbol(ctx context.Context, symbol string) float64
	GetLowestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64
	GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64

	GetAvgPriceBySymbol(ctx context.Context, symbol string) float64
	GetAvgPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64
	GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64
}

type DataProcessingService interface {
	StartWorkers(ctx context.Context, in <-chan PriceTick)
	StopWorkers()
}

type ServerHandler interface {
	StartServer(ctx context.Context, config *config.Config) error
	RegisterRouter(mux *http.ServeMux)
}

type MarketDataService interface {
	GetLatestPriceBySymbol(ctx context.Context, symbol string) float64
	GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64

	GetHighestPriceBySymbol(ctx context.Context, symbol string) float64
	GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64
	GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64

	GetLowestPriceBySymbol(ctx context.Context, symbol string) float64
	GetLowestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64
	GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64

	GetAvgPriceBySymbol(ctx context.Context, symbol string) float64
	GetAvgPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64
	GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64
}

type SystemService interface {
	SwitchToTestMode(ctx context.Context) error
	SwitchToLiveMode(ctx context.Context) error
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
