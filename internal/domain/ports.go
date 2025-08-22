package domain

import (
	"context"
	"time"
)

type ExchangeClient interface {
	StartLiveMode(ctx context.Context) <-chan PriceTick
	StartTestMode(ctx context.Context) <-chan PriceTick
	Stop()
}

type RedisClient interface {
	StoreTick(ctx context.Context, exchange, pair string, price float64)
	ProcessLastMinute(ctx context.Context, exchange, pair string) *MinuteAgg

	GetLatestPriceBySymbol(ctx context.Context, symbol string) float64
	GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64

	GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64

	GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64

	GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64
	GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64
}

type RepositoryClient interface {
	StoreMinAgg(ctx context.Context, agg *MinuteAgg)

	GetHighestPriceBySymbol(ctx context.Context, symbol string)
	GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string)
	GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration)
	GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration)

	GetLowestPriceBySymbol(ctx context.Context, symbol string)
	GetLowestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string)
	GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration)
	GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration)

	GetAvgPriceBySymbol(ctx context.Context, symbol string)
	GetAvgPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string)
	GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration)
	GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration)
}

type DataProcessingService interface {
	StartWorkers(ctx context.Context, in <-chan PriceTick)
	StopWorkers()
}

type MarketDataService interface {
	GetLatestPriceBySymbol()
	GetLatestPriceBySymbolAndExchange()

	GetHighestPriceBySymbol()
	GetHighestPriceBySymbolAndExchange()
	GetHighestPriceBySymbolAndPeriod()
	GetHighestPriceBySymbolAndPeriodAndExchange()

	GetLowestPriceBySymbol()
	GetLowestPriceBySymbolAndExchange()
	GetLowestPriceBySymbolAndPeriod()
	GetLowestPriceBySymbolAndPeriodAndExchange()

	GetAvgPriceBySymbol()
	GetAvgPriceBySymbolAndExchange()
	GetAvgPriceBySymbolAndPeriod()
	GetAvgPriceBySymbolAndPeriodAndExchange()
}

type SystemService interface {
	SwitchToLiveMode()
	SwitchToTestMode()

	GetHealth()
}
