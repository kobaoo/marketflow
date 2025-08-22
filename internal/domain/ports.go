package domain

import (
	"context"
)

type ExchangeClient interface {
	StartLiveMode(ctx context.Context) <-chan PriceTick
	StartTestMode(ctx context.Context) <-chan PriceTick
	Stop()
}

type RedisClient interface {
	StoreTick(ctx context.Context, exchange, pair string, price float64)
	ProcessLastMinute(ctx context.Context, pair, exchange string) *MinuteAgg

	GetLatestPriceBySymbol()
	GetLatestPriceBySymbolAndExchange()

	GetHighestPriceBySymbolAndPeriod()
	GetHighestPriceBySymbolAndPeriodAndExchange()

	GetLowestPriceBySymbolAndPeriod()
	GetLowestPriceBySymbolAndPeriodAndExchange()

	GetAvgPriceBySymbolAndPeriod()
	GetAvgPriceBySymbolAndPeriodAndExchange()
}

type Repository interface {
	StoreMinTick(ctx context.Context, tick *MinuteAgg)

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

type SystemService interface{
	SwitchToLiveMode()
	SwitchToTestMode()

	GetHealth()
}
