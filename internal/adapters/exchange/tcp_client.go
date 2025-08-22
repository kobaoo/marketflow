package exchange

import (
	"bufio"
	"context"
	"log/slog"
	"marketflow/internal/adapters/cache"
	"marketflow/internal/config"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// TODO: implement auto reconnecting to server if connections are lost

func RunTCPClients(config *config.Config, rdb *cache.RedisClient, testMode bool) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	allMessages := make(chan []byte, 30)

	if testMode {
		genMessages := make(chan []byte, 30)
		go StartGenerator(ctx, genMessages)
		go StartGenerator(ctx, genMessages)
		go StartGenerator(ctx, genMessages)

	} else {
		runTCPClient(ctx, "127.0.0.1:40101", 1, allMessages)
		runTCPClient(ctx, "127.0.0.1:40102", 2, allMessages)
		runTCPClient(ctx, "127.0.0.1:40103", 3, allMessages)
	}
}

func runTCPClient(ctx context.Context, address string, exchangeId int, allMessages chan<- []byte) {
	// Connect to TCP server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("Connection error", "error", err)
		return
	}
	defer conn.Close()

	// Channel with small buffer to avoid lag
	messages := make(chan []byte, 10)

	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(ctx, i, exchangeId, messages, allMessages, &wg)
	}

	// TCP reader goroutine with drop-on-overflow
	go func() {
		reader := bufio.NewReader(conn)
		for {
			select {
			case <-ctx.Done():
				close(messages)
				return
			default:
				line, err := reader.ReadBytes('\n')
				if err != nil {
					slog.Error("Read error:", "err", err)
					close(messages)
					return
				}
				// Try to send, drop if channel is full
				select {
				case messages <- line:
					// sent successfully
				default:
					slog.Debug("⚠ Dropping stale message")
				}
			}
		}
	}()

	// Wait for shutdown
	<-ctx.Done()
	slog.Info("Shutting down...")
	wg.Wait()
}

func worker(ctx context.Context, id, exchangeId int, jobs <-chan []byte, allMessages chan<- []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			slog.Info("Worker stopping", "exchange", exchangeId, "id", id)
			return
		case msg, ok := <-jobs:
			if !ok {
				return
			}
			select {
			case allMessages <- msg:
				slog.Debug("Worker processing message", "exchange", exchangeId, "id", id, "message", string(msg))
			default:
				slog.Debug("Worker dropping message", "exchange", exchangeId, "id", id, "message", string(msg))
			}
		}
	}
}
