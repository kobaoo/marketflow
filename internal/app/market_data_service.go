package app

import (
	"context"
	"marketflow/internal/domain"
)

type MarketDataService struct {
	redisClient      domain.RedisClient
	repositoryClient domain.Repository
}

func NewMarketDataService(redisClient domain.RedisClient, repositoryClient domain.Repository) domain.MarketDataService {
	return &MarketDataService{
		redisClient:      redisClient,
		repositoryClient: repositoryClient,
	}
}

func (r *MarketDataService) GetLatestPriceBySymbol(ctx context.Context, symbol string) float64 {
	return r.redisClient.GetLatestPriceBySymbol(ctx, symbol)
}
func (r *MarketDataService) GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	return r.redisClient.GetLatestPriceBySymbolAndExchange(ctx, symbol, exchange)
}

func (r *MarketDataService) GetHighestPriceBySymbol() float64                     { return 0 }
func (r *MarketDataService) GetHighestPriceBySymbolAndExchange() float64          { return 0 }
func (r *MarketDataService) GetHighestPriceBySymbolAndPeriod() float64            { return 0 }
func (r *MarketDataService) GetHighestPriceBySymbolAndPeriodAndExchange() float64 { return 0 }

func (r *MarketDataService) GetLowestPriceBySymbol() float64                     { return 0 }
func (r *MarketDataService) GetLowestPriceBySymbolAndExchange() float64          { return 0 }
func (r *MarketDataService) GetLowestPriceBySymbolAndPeriod() float64            { return 0 }
func (r *MarketDataService) GetLowestPriceBySymbolAndPeriodAndExchange() float64 { return 0 }

func (r *MarketDataService) GetAvgPriceBySymbol() float64                     { return 0 }
func (r *MarketDataService) GetAvgPriceBySymbolAndExchange() float64          { return 0 }
func (r *MarketDataService) GetAvgPriceBySymbolAndPeriod() float64            { return 0 }
func (r *MarketDataService) GetAvgPriceBySymbolAndPeriodAndExchange() float64 { return 0 }
