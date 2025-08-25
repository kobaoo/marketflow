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
	config *config.Config
	mu sync.Mutex
	retryDelay time.Duration
}

func NewExchangeClient(config *config.Config, retryDelay time.Duration) domain.ExchangeClient {
	return &ExchangeClient{config: config, retryDelay: retryDelay}
}

func (r *ExchangeClient) StartLiveMode(ctx context.Context) <-chan domain.PriceTick {
	messages := make(chan domain.PriceTick, 30)

	// Start TCP clients for each exchange in separate goroutines
	for _, exchange := range r.config.Exchanges {
		go r.runTCPClient(ctx, exchange.Addr, exchange.Name, messages)
	}

	// Close channel when context is done
	go func() {
		<-ctx.Done()
		close(messages)
	}()

	return messages
}

func (r *ExchangeClient) StartTestMode(ctx context.Context) <-chan domain.PriceTick {
	messages := make(chan domain.PriceTick, 30)

	// Start generators for test exchanges in separate goroutines
	go r.startGenerator(ctx, "ex1", messages)
	go r.startGenerator(ctx, "ex2", messages)
	go r.startGenerator(ctx, "ex3", messages)

	// Close channel when context is done
	go func() {
		<-ctx.Done()
		close(messages)
	}()

	return messages
}

func (r *ExchangeClient) Stop(stop context.CancelFunc) {
	stop()
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

			if err := r.handleConnection(ctx, conn, exchangeName, out); err != nil {
				slog.Warn("Connection handler error",
					"exchange", exchangeName,
					"error", err)
			}

			conn.Close()
			slog.Info("Waiting before reconnect", "exchange", exchangeName)
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
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return err
			}

			var tick domain.PriceTick
			if err := json.Unmarshal(line, &tick); err != nil {
				slog.Error("Unmarshal error",
					"exchange", exchangeName,
					"error", err,
					"data", string(line))
				continue
			}

			tick.Exchange = exchangeName

			select {
			case out <- tick:
				// send successfully
			case <-ctx.Done():
				return nil
			default:
				slog.Warn("Channel full, dropping message",
					"exchange", exchangeName,
					"symbol", tick.Symbol)
			}
		}
	}
}

