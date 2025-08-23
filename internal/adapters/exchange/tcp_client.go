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

// TODO: implement auto reconnecting to server if connections are lost

func RunTCPClients(allMessages chan<- []byte, testMode bool) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	var wg sync.WaitGroup

	if testMode {
		genMessages := make(chan []byte, 30)
		go StartGenerator(ctx, genMessages)
		go StartGenerator(ctx, genMessages)
		go StartGenerator(ctx, genMessages)

	} else {
		wg.Add(3)
		go func() {
			runTCPClient(ctx, "exchange1:40101", 1, allMessages)
			wg.Done()
		}()
		go func() {
			runTCPClient(ctx, "exchange2:40102", 2, allMessages)
			wg.Done()
		}()
		go func() {
			runTCPClient(ctx, "exchange3:40103", 3, allMessages)
			wg.Done()
		}()
	}
	
	wg.Wait()
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
	slog.Info("Shutting down...", "exchange", exchangeId)
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
				slog.Debug("Channel closed", "chanel", "messages")
				return
			}
			select {
			case allMessages <- msg:
				slog.Debug("Worker processing message", "exchange", exchangeId, "id", id, "message", msg)
			default:
				slog.Debug("Worker dropping message", "exchange", exchangeId, "id", id, "message", string(msg))
			}
		}
	}
}
