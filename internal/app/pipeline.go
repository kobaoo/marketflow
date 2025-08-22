package app

import (
	"context"
	"log/slog"
	"marketflow/internal/domain"
	"sync"
	"time"
)

type DataProcessingService struct {
	redisClient domain.RedisClient
	repository  domain.Repository
}

func NewDataProcessingService(redisClient domain.RedisClient, repository domain.Repository) domain.DataProcessingService {
	return &DataProcessingService{
		redisClient: redisClient,
		repository:  repository,
	}
}

func (r *DataProcessingService) StartWorkers(ctx context.Context, in <-chan domain.PriceTick) {
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= 15; i++ {
		wg.Add(1)
		go r.worker(ctx, in, &wg)
	}

	<-ctx.Done()
	slog.Info("Shutting down...")
	wg.Wait()
}
func (r *DataProcessingService) StopWorkers(stop context.CancelFunc) {
	stop()
}

func (r *DataProcessingService) worker(ctx context.Context, in <-chan domain.PriceTick, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Worker stopping")
			return
		case <-ticker.C:
			agg := r.redisClient.ProcessLastMinute(ctx)
			r.repository.StoreMinAgg(ctx, agg)
		case msg, ok := <-in:
			if !ok {
				slog.Error("Error channel closed")
				return
			}
			r.redisClient.StoreTick(ctx, msg.Exchange, msg.Symbol, msg.Price)
		}
	}
}
