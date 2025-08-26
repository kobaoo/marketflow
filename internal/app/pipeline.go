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
	wg          sync.WaitGroup
	cancelFunc  context.CancelFunc
}

func NewDataProcessingService(redisClient domain.RedisClient, repository domain.Repository) domain.DataProcessingService {
	return &DataProcessingService{
		redisClient: redisClient,
		repository:  repository,
	}
}

func (d *DataProcessingService) StopWorkers() {
	if d.cancelFunc != nil {
		d.cancelFunc()
	}

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("All data processing workers stopped")
	case <-time.After(5 * time.Second):
		slog.Warn("Timeout waiting data processing workers to stop")
	}
}

func (d *DataProcessingService) StartWorkers(ctx context.Context, in <-chan domain.PriceTick) {
	if ctx.Err() != nil {
		slog.Warn("Cannot start workers: context already cancelled")
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	d.cancelFunc = cancel

	for i := 1; i <= 15; i++ {
		d.wg.Add(1)
		go d.worker(ctx, in, i)
	}

	d.wg.Add(1)
	go d.aggregator(ctx)

	slog.Info("Started data processing workers", "count", 15)
}

func (d *DataProcessingService) worker(ctx context.Context, in <-chan domain.PriceTick, workerID int) {
	defer d.wg.Done()

	slog.Info("Worker started", "worker", workerID)
	defer slog.Info("Worker stopped", "worker", workerID)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-in:
			if !ok {
				slog.Info("Input channel closed", "worker", workerID)
				return
			}
			d.redisClient.StoreTick(ctx, msg.Exchange, msg.Symbol, msg.Price)
		}
	}
}

func (d *DataProcessingService) aggregator(ctx context.Context) {
	defer d.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if agg := d.redisClient.ProcessLastMinute(ctx); len(agg) > 0 {
				d.repository.StoreMinAgg(ctx, agg)
			}
		}
	}
}
