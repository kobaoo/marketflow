package app

import (
	"context"
	"marketflow/internal/domain"
	"time"
)

type MarketDataService struct {
	redisClient domain.RedisClient
	repository  domain.Repository
}

func NewMarketDataService(redisClient domain.RedisClient, repository domain.Repository) domain.MarketDataService {
	return &MarketDataService{
		redisClient: redisClient,
		repository:  repository,
	}
}

func (r *MarketDataService) GetLatestPriceBySymbol(ctx context.Context, symbol string) float64 {
	return r.redisClient.GetLatestPriceBySymbol(ctx, symbol)
}
func (r *MarketDataService) GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	return r.redisClient.GetLatestPriceBySymbolAndExchange(ctx, symbol, exchange)
}

func (r *MarketDataService) GetHighestPriceBySymbol(ctx context.Context, symbol string) float64 {
	return r.repository.GetHighestPriceBySymbol(ctx, symbol)
}
func (r *MarketDataService) GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	return r.repository.GetHighestPriceBySymbolAndExchange(ctx, symbol, exchange)
}
func (r *MarketDataService) GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	if period <= time.Minute {
		return r.redisClient.GetHighestPriceBySymbolAndPeriod(ctx, symbol, period)
	} else {
		return r.repository.GetHighestPriceBySymbolAndPeriod(ctx, symbol, period)
	}
}
func (r *MarketDataService) GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	if period <= time.Minute {
		return r.redisClient.GetHighestPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
	} else {
		return r.repository.GetHighestPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
	}
}

func (r *MarketDataService) GetLowestPriceBySymbol(ctx context.Context, symbol string) float64 {
	return r.repository.GetLowestPriceBySymbol(ctx, symbol)
}
func (r *MarketDataService) GetLowestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	return r.repository.GetLowestPriceBySymbolAndExchange(ctx, symbol, exchange)
}
func (r *MarketDataService) GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	if period <= time.Minute {
		return r.redisClient.GetLowestPriceBySymbolAndPeriod(ctx, symbol, period)
	} else {
		return r.repository.GetLowestPriceBySymbolAndPeriod(ctx, symbol, period)
	}
}
func (r *MarketDataService) GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	if period <= time.Minute {
		return r.redisClient.GetLowestPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
	} else {
		return r.repository.GetLowestPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
	}
}

func (r *MarketDataService) GetAvgPriceBySymbol(ctx context.Context, symbol string) float64 {
	return r.repository.GetAvgPriceBySymbol(ctx, symbol)
}
func (r *MarketDataService) GetAvgPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	return r.repository.GetAvgPriceBySymbolAndExchange(ctx, symbol, exchange)
}
func (r *MarketDataService) GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	if period <= time.Minute {
		return r.redisClient.GetAvgPriceBySymbolAndPeriod(ctx, symbol, period)
	} else {
		return r.repository.GetAvgPriceBySymbolAndPeriod(ctx, symbol, period)
	}
}
func (r *MarketDataService) GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	if period <= time.Minute {
		return r.redisClient.GetAvgPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
	} else {
		return r.repository.GetAvgPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
	}
}
