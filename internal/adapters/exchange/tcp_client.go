package exchange

import (
	"bufio"
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func RunTCPClient() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to TCP server
	conn, err := net.Dial("tcp", "127.0.0.1:40101")
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
		go worker(ctx, i, messages, &wg)
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

func worker(ctx context.Context, id int, jobs <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			slog.Info("Worker stopping", "id", id)
			return
		case msg, ok := <-jobs:
			if !ok {
				return
			}
			slog.Debug("Worker processing message", "id", id, "message", string(msg))
		}
	}
}
