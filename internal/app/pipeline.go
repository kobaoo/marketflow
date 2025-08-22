package app

import "marketflow/internal/domain"

type DataProcessingService struct {
	exchangeClient   domain.ExchangeClient
	redisClient      domain.RedisClient
	repositoryClient domain.Repository
}

func NewDataProcessingService(exchangeClient domain.ExchangeClient, redisClient domain.RedisClient, repositoryClient domain.Repository) *DataProcessingService {
	return &DataProcessingService{
		exchangeClient:   exchangeClient,
		redisClient:      redisClient,
		repositoryClient: repositoryClient,
	}
}
