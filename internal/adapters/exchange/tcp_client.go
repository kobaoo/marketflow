package exchange

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"net"
	"strconv"
)

type ExchangeClient struct {
	config *config.Config
}

func NewExchangeClient(config *config.Config) domain.ExchangeClient {
	return &ExchangeClient{config: config}
}

func (r ExchangeClient) StartLiveMode(ctx context.Context) <-chan domain.PriceTick {
	messages := make(chan domain.PriceTick, 30)

	r.runTCPClient(ctx, "127.0.0.1:40101", 1, messages)
	r.runTCPClient(ctx, "127.0.0.1:40102", 2, messages)
	r.runTCPClient(ctx, "127.0.0.1:40103", 3, messages)

	return messages
}

func (r ExchangeClient) StartTestMode(ctx context.Context) <-chan domain.PriceTick {
	messages := make(chan domain.PriceTick, 30)

	r.startGenerator(ctx, "exchange1", messages)
	r.startGenerator(ctx, "exchange2", messages)
	r.startGenerator(ctx, "exchange3", messages)

	return messages
}

func (r ExchangeClient) Stop(stop context.CancelFunc) {
	stop()
}

func (r ExchangeClient) runTCPClient(ctx context.Context, address string, exchangeId int, out chan<- domain.PriceTick) {
	// Connect to TCP server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("Connection error", "error", err)
		return
	}
	defer conn.Close()

	// TCP reader goroutine with drop-on-overflow
	go func() {
		reader := bufio.NewReader(conn)
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			default:
				line, err := reader.ReadBytes('\n')
				if err != nil {
					slog.Error("Read error:", "err", err)
					close(out)
					return
				}

				tick := domain.PriceTick{Exchange: "exchange" + strconv.Itoa(exchangeId)}
				err = json.Unmarshal(line, &tick)
				if err != nil {
					slog.Error("Error unmarshaling message")
				}

				// Try to send, drop if channel is full
				select {
				case out <- tick:
					// sent successfully
				default:
					slog.Debug("⚠ Dropping stale message")
				}
			}
		}
	}()

	// Wait for shutdown
	<-ctx.Done()
	slog.Info("Shutting down tcp client")
}
