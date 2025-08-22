package app

import "marketflow/internal/domain"

type DataProcessingService struct {
	exchangeClient domain.ExchangeClient
	repositoryClient domain.RepositoryClient
}