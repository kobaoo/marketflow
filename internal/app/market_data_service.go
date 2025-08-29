package app

import (
	"context"
	"marketflow/internal/domain"
	"strings"
	"time"
)

type MarketDataService struct {
	redisClient domain.RedisClient
	repository  domain.Repository
	windowStore domain.WindowStore
}

func NewMarketDataService(redisClient domain.RedisClient, repository domain.Repository, window domain.WindowStore) domain.MarketDataService {
	return &MarketDataService{
		redisClient: redisClient,
		repository:  repository,
		windowStore: window,
	}
}

func (r *MarketDataService) GetLatestPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	symbol = strings.ToUpper(symbol)
	if price, err := r.redisClient.GetLatestPriceBySymbol(ctx, symbol); err == nil {
		return price, nil
	}
	return r.windowStore.GetLatestPriceBySymbol(symbol)
}

func (r *MarketDataService) GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error) {
	symbol = strings.ToUpper(symbol)
	if price, err := r.redisClient.GetLatestPriceBySymbolAndExchange(ctx, symbol, exchange); err == nil {
		return price, nil
	}
	return r.windowStore.GetLatestPriceBySymbolAndExchange(symbol, exchange)
}

func (r *MarketDataService) GetHighestPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	symbol = strings.ToUpper(symbol)
	return r.repository.GetHighestPriceBySymbol(ctx, symbol)
}

func (r *MarketDataService) GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error) {
	symbol = strings.ToUpper(symbol)
	return r.repository.GetHighestPriceBySymbolAndExchange(ctx, symbol, exchange)
}

func (r *MarketDataService) GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) (float64, error) {
	symbol = strings.ToUpper(symbol)

	if period <= time.Minute {
		if price, err := r.redisClient.GetHighestPriceBySymbolAndPeriod(ctx, symbol, period); err == nil {
			return price, nil
		}
		return r.windowStore.GetHighestPriceBySymbolAndPeriod(symbol, period)
	}
	return r.repository.GetHighestPriceBySymbolAndPeriod(ctx, symbol, period)
}

func (r *MarketDataService) GetHighestPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	symbol = strings.ToUpper(symbol)

	if period <= time.Minute {
		if price, err := r.redisClient.GetHighestPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period); err == nil {
			return price, nil
		}
		return r.windowStore.GetHighestPriceBySymbolAndPeriodAndExchange(symbol, exchange, period)
	}
	return r.repository.GetHighestPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
}

func (r *MarketDataService) GetLowestPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	symbol = strings.ToUpper(symbol)
	return r.repository.GetLowestPriceBySymbol(ctx, symbol)
}

func (r *MarketDataService) GetLowestPriceBySymbolAndExchange(
	ctx context.Context, symbol, exchange string,
) (float64, error) {
	symbol = strings.ToUpper(symbol)
	return r.repository.GetLowestPriceBySymbolAndExchange(ctx, symbol, exchange)
}

func (r *MarketDataService) GetLowestPriceBySymbolAndPeriod(
	ctx context.Context, symbol string, period time.Duration,
) (float64, error) {
	symbol = strings.ToUpper(symbol)

	if period <= time.Minute {
		if price, err := r.redisClient.GetLowestPriceBySymbolAndPeriod(ctx, symbol, period); err == nil {
			return price, nil
		}
		return r.windowStore.GetLowestPriceBySymbolAndPeriod(symbol, period)
	}
	return r.repository.GetLowestPriceBySymbolAndPeriod(ctx, symbol, period)
}

func (r *MarketDataService) GetLowestPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	symbol = strings.ToUpper(symbol)

	if period <= time.Minute {
		if price, err := r.redisClient.GetLowestPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period); err == nil {
			return price, nil
		}
		return r.windowStore.GetLowestPriceBySymbolAndPeriodAndExchange(symbol, exchange, period)
	}
	return r.repository.GetLowestPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
}

func (r *MarketDataService) GetAvgPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	symbol = strings.ToUpper(symbol)
	return r.repository.GetAvgPriceBySymbol(ctx, symbol)
}

func (r *MarketDataService) GetAvgPriceBySymbolAndExchange(
	ctx context.Context, symbol, exchange string,
) (float64, error) {
	symbol = strings.ToUpper(symbol)
	return r.repository.GetAvgPriceBySymbolAndExchange(ctx, symbol, exchange)
}

func (r *MarketDataService) GetAvgPriceBySymbolAndPeriod(
	ctx context.Context, symbol string, period time.Duration,
) (float64, error) {
	symbol = strings.ToUpper(symbol)

	if period <= time.Minute {
		if price, err := r.redisClient.GetAvgPriceBySymbolAndPeriod(ctx, symbol, period); err == nil {
			return price, nil
		}
		return r.windowStore.GetAvgPriceBySymbolAndPeriod(symbol, period)
	}
	return r.repository.GetAvgPriceBySymbolAndPeriod(ctx, symbol, period)
}

func (r *MarketDataService) GetAvgPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	symbol = strings.ToUpper(symbol)

	if period <= time.Minute {
		if price, err := r.redisClient.GetAvgPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period); err == nil {
			return price, nil
		}
		return r.windowStore.GetAvgPriceBySymbolAndPeriodAndExchange(symbol, exchange, period)
	}
	return r.repository.GetAvgPriceBySymbolAndPeriodAndExchange(ctx, symbol, exchange, period)
}
