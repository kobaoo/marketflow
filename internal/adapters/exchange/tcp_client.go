package exchange

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"net"
	"sync"
	"time"
)

type ExchangeClient struct {
	config     *config.Config
	mu         sync.Mutex
	retryDelay time.Duration
}

func NewExchangeClient(config *config.Config, retryDelay time.Duration) domain.ExchangeClient {
	return &ExchangeClient{
		config:     config,
		retryDelay: retryDelay,
	}
}

func (r *ExchangeClient) StartLiveMode(ctx context.Context) <-chan domain.PriceTick {
	messages := make(chan domain.PriceTick, 100)

	for _, exchange := range r.config.Exchanges {
		go r.runTCPClient(ctx, exchange.Addr, exchange.Name, messages)
	}

	return messages
}

func (r *ExchangeClient) StartTestMode(ctx context.Context) <-chan domain.PriceTick {
	messages := make(chan domain.PriceTick, 30)

	if ctx.Err() != nil {
		slog.Warn("Context already cancelled, cannot start test mode")
		close(messages)
		return messages
	}

	go r.startGenerator(ctx, "ex1", messages)
	go r.startGenerator(ctx, "ex2", messages)
	go r.startGenerator(ctx, "ex3", messages)

	go func() {
		<-ctx.Done()
		slog.Info("Test mode context cancelled, closing messages channel")
		close(messages)
	}()

	return messages
}

func (r *ExchangeClient) Stop() {
	slog.Info("Exchange client stop requested")
}

func (r *ExchangeClient) runTCPClient(ctx context.Context, address string, exchangeName string, out chan<- domain.PriceTick) {
	slog.Info("Starting TCP client", "exchange", exchangeName, "address", address)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down TCP client", "exchange", exchangeName)
			return
		default:
			conn, err := net.Dial("tcp", address)
			if err != nil {
				slog.Error("Connection error, retrying...",
					"exchange", exchangeName,
					"address", address,
					"error", err)
				time.Sleep(r.retryDelay)
				continue
			}

			slog.Info("TCP client connected", "exchange", exchangeName, "address", address)

			// Обрабатываем соединение
			if err := r.handleConnection(ctx, conn, exchangeName, out); err != nil {
				slog.Warn("Connection handler error", "exchange", exchangeName, "error", err)
			}

			conn.Close()
			time.Sleep(r.retryDelay)
		}
	}
}

func (r *ExchangeClient) handleConnection(ctx context.Context, conn net.Conn, exchangeName string, out chan<- domain.PriceTick) error {
	reader := bufio.NewReader(conn)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Context cancelled, closing connection", "exchange", exchangeName)
			return nil
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		var raw rawTick
		if err := json.Unmarshal(line, &raw); err != nil {
			slog.Error("Unmarshal error", "exchange", exchangeName, "error", err)
			continue
		}

		ts := time.Unix(0, raw.Timestamp*int64(time.Millisecond))

		tick := domain.PriceTick{
			Exchange: exchangeName,
			Symbol:   raw.Symbol,
			Price:    raw.Price,
			Ts:       ts,
		}

		select {
		case <-ctx.Done():
			return nil
		case out <- tick:
		case <-time.After(100 * time.Millisecond):
			slog.Warn("Send timeout, dropping message", "exchange", exchangeName)
		}
	}
}
