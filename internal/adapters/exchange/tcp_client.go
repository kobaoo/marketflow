package exchange

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"net"
)

type ExchangeClient struct {
	config *config.Config
}

func NewExchangeClient(config *config.Config) domain.ExchangeClient {
	return &ExchangeClient{config: config}
}

func (r ExchangeClient) StartLiveMode(ctx context.Context) <-chan domain.PriceTick {
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

func (r ExchangeClient) StartTestMode(ctx context.Context) <-chan domain.PriceTick {
	messages := make(chan domain.PriceTick, 30)

	// Start generators for test exchanges in separate goroutines
	go r.startGenerator(ctx, "exchange1", messages)
	go r.startGenerator(ctx, "exchange2", messages)
	go r.startGenerator(ctx, "exchange3", messages)

	// Close channel when context is done
	go func() {
		<-ctx.Done()
		close(messages)
	}()

	return messages
}

func (r ExchangeClient) Stop(stop context.CancelFunc) {
	stop()
}

func (r ExchangeClient) runTCPClient(ctx context.Context, address string, exchangeName string, out chan<- domain.PriceTick) {
	slog.Info("Starting TCP client", "exchange", exchangeName, "address", address)
	
	// Connect to TCP server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("Connection error", "exchange", exchangeName, "address", address, "error", err)
		return
	}
	defer conn.Close()

	slog.Info("TCP client connected", "exchange", exchangeName, "address", address)

	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down tcp client", "exchange", exchangeName)
			return
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				slog.Error("Read error", "exchange", exchangeName, "error", err)
				return
			}

			tick := domain.PriceTick{Exchange: exchangeName}
			err = json.Unmarshal(line, &tick)
			if err != nil {
				slog.Error("Error unmarshaling message", "exchange", exchangeName, "error", err, "data", string(line))
				continue
			}

			// Try to send, drop if channel is full
			select {
			case out <- tick:
				// sent successfully
			default:
				slog.Debug("⚠ Dropping stale message", "exchange", exchangeName, "symbol", tick.Symbol)
			}
		}
	}
}
