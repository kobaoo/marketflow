package app

import (
	"context"
	"log/slog"
	"marketflow/internal/domain"
	"sync"
	"time"
)

type DataProcessingService struct {
	redisClient  domain.RedisClient
	repository   domain.Repository
	wg           sync.WaitGroup
	cancelFunc   context.CancelFunc
	lastTickMu   sync.RWMutex
	lastTick     map[string]time.Time
	window       domain.WindowStore
	redisTimeout time.Duration
	breakerMu    sync.Mutex
	redisBlocked bool
	unblockAfter time.Time
}

func NewDataProcessingService(redisClient domain.RedisClient, repository domain.Repository, window domain.WindowStore) domain.DataProcessingService {
	return &DataProcessingService{
		redisClient:  redisClient,
		repository:   repository,
		window:       window,
		redisTimeout: 150 * time.Millisecond,
		lastTick:     make(map[string]time.Time),
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

			d.window.Add(msg)
			d.tryStoreToRedis(ctx, msg)

			d.lastTickMu.Lock()
			d.lastTick[msg.Exchange] = msg.Ts
			d.lastTickMu.Unlock()
		}
	}
}

func (d *DataProcessingService) tryStoreToRedis(parent context.Context, msg domain.PriceTick) {
	if d.redisClient == nil {
		return
	}

	d.breakerMu.Lock()
	if d.redisBlocked && time.Now().Before(d.unblockAfter) {
		d.breakerMu.Unlock()
		return
	}
	d.breakerMu.Unlock()

	ctx, cancel := context.WithTimeout(parent, d.redisTimeout)
	defer cancel()

	if err := d.redisClient.StoreTick(ctx, msg.Exchange, msg.Symbol, msg.Price); err != nil {
		slog.Warn("Redis write failed", "err", err)

		d.breakerMu.Lock()
		d.redisBlocked = true
		d.unblockAfter = time.Now().Add(5 * time.Second)
		d.breakerMu.Unlock()
		return
	}

	d.breakerMu.Lock()
	d.redisBlocked = false
	d.breakerMu.Unlock()
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
			snap := d.window.Snapshot()
			if len(snap) == 0 {
				continue
			}

			var rows []*domain.MinuteAgg
			now := time.Now()

			for k, arr := range snap {
				if len(arr) == 0 {
					continue
				}

				var sum float64
				minV, maxV := arr[0].Price, arr[0].Price
				for _, t := range arr {
					sum += t.Price
					if t.Price < minV {
						minV = t.Price
					}
					if t.Price > maxV {
						maxV = t.Price
					}
				}
				avg := sum / float64(len(arr))

				rows = append(rows, &domain.MinuteAgg{
					Symbol:   k.Symbol,
					Exchange: k.Exchange,
					Ts:       now.Truncate(time.Minute),
					Avg:      avg,
					Min:      minV,
					Max:      maxV,
				})
			}

			if len(rows) > 0 {
				if err := d.repository.StoreMinAgg(ctx, rows); err != nil {
					slog.Error("PG insert failed", "err", err)
				}
			}
		}
	}
}

func (d *DataProcessingService) ExchangesHealth(within time.Duration, names []string) map[string]bool {
	now := time.Now()
	out := make(map[string]bool, len(names))

	d.lastTickMu.RLock()
	defer d.lastTickMu.RUnlock()

	for _, name := range names {
		if t, ok := d.lastTick[name]; ok && now.Sub(t) <= within {
			out[name] = true
		} else {
			out[name] = false
		}
	}
	return out
}
