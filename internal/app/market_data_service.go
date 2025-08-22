package app

import "marketflow/internal/domain"

type MarketDataService struct {
	redisClient      domain.RedisClient
	repositoryClient domain.Repository
}

func NewMarketDataService(redisClient domain.RedisClient, repositoryClient domain.Repository) *MarketDataService {
	return &MarketDataService{
		redisClient:      redisClient,
		repositoryClient: repositoryClient,
	}
}
