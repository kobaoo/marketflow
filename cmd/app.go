package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"marketflow/internal/adapters/cache"
	"marketflow/internal/adapters/exchange"
	"marketflow/internal/adapters/postgres"
	"marketflow/internal/adapters/web"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"marketflow/internal/infra"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func RunApp() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	infra.SetUpLogger()

	config, err := config.ReadConfig() // Read config
	if err != nil {
		slog.Error("Config Error", "error", err)
	}
	priceCh := make(chan *domain.PriceTick, 100)
	defer close(priceCh)

	var streams []domain.ExchangeStream
	for _, exchangeCfg := range config.Exchanges {
		stream := exchange.NewTCPHandler(
			exchangeCfg,
			priceCh,
			5*time.Second,
		)
		streams = append(streams, stream)
	}

	rdb, err := cache.NewRedisClient(&config) // Connect to Redis
	if err != nil {
		slog.Error("Redis Error", "error", err)
	}

	db := postgres.ConnectDB(&config)
	defer db.Close()
	repository := postgres.NewRepository(db)
	fmt.Println(repository, rdb)

	var wg sync.WaitGroup
	for _, stream := range streams {
		wg.Add(1)
		go func(s domain.ExchangeStream) {
			defer wg.Done()
			slog.Info("Starting exchange stream", "exchange", s.GetExchangeName())
			s.Start(ctx)
		}(stream)
	}

	// TODO: Processor to accept and aggregate data coming from priceCH channel

	wg.Add(1)
	go func() {
		defer wg.Done()
		processPriceTicks(ctx, priceCh)
	}()


	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("Starting web server", "port", config.Port)
		if err := web.StartServer(&config); err != nil {
			slog.Error("Web server error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down application...")

	for _, stream := range streams {
		stream.Stop()
	}

	wg.Wait()
	slog.Info("Application stopped gracefully")
}

// it's just for checking tcpClient
func processPriceTicks(
	ctx context.Context, 
	priceCh <-chan *domain.PriceTick,
) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Price tick processor stopped")
			return
		case tick, ok := <-priceCh:
			if !ok {
				return
			}
			
			slog.Debug("Price tick processed", 
				"exchange", tick.Exchange, 
				"symbol", tick.Symbol,
				"price", tick.Price,
				"stack of channel", len(priceCh))
		}
	}
}