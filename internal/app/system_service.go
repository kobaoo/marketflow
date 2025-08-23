package app

import "marketflow/internal/domain"

type SystemService struct {
	redisClient      domain.RedisClient
	repositoryClient domain.Repository
}

func NewSystemService(redisClient domain.RedisClient, repositoryClient domain.Repository) *SystemService {
	return &SystemService{redisClient: redisClient, repositoryClient: repositoryClient}
}
