package exchange

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"net"
	"sync"
	"time"
)

type TCPHandler struct {
	exchange   config.Exchange
	outputCh   chan<- *domain.PriceTick
	retryDelay time.Duration
	mu         sync.Mutex
	running    bool
	cancel     context.CancelFunc
}

func NewTCPHandler(
	exchange config.Exchange,
	outputCh chan<- *domain.PriceTick,
	retryDelay time.Duration,
) domain.ExchangeStream {
	return &TCPHandler{
		exchange:   exchange,
		outputCh:   outputCh,
		retryDelay: retryDelay,
		running:    false,
	}
}

func (h *TCPHandler) Start(ctx context.Context) {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return
	}
	
	childCtx, cancel := context.WithCancel(ctx)
	h.cancel = cancel
	h.running = true
	h.mu.Unlock()

	go h.run(childCtx)
}

func (h *TCPHandler) run(ctx context.Context) {
	defer func() {
		h.mu.Lock()
		h.running = false
		h.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			slog.Info("TCP handler stopped", "exchange", h.exchange.Name)
			return
		default:
			if err := h.connectAndProcess(ctx); err != nil {
				slog.Error("TCP connection failed", 
					"exchange", h.exchange.Name, 
					"error", err,
					"retry_in", h.retryDelay)
				time.Sleep(h.retryDelay)
			}
		}
	}
}

func (h *TCPHandler) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if h.cancel != nil {
		h.cancel()
	}
	h.running = false
}

func (h *TCPHandler) GetExchangeName() string {
	return h.exchange.Name
}

func (h *TCPHandler) connectAndProcess(ctx context.Context) error {
	conn, err := net.Dial("tcp", h.exchange.Addr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()
	
	slog.Info("TCP connected", "exchange", h.exchange.Name, "address", h.exchange.Addr)
	
	reader := bufio.NewReader(conn)
	
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return fmt.Errorf("read failed: %w", err)
			}
			
			price, err := h.parsePriceData(line)
			if err != nil {
				slog.Warn("Failed to parse price data", 
					"exchange", h.exchange.Name, 
					"error", err,
					"data", string(line))
				continue
			}
			
			select {
			case h.outputCh <- price:
				slog.Debug("Price processed", 
					"exchange", h.exchange.Name, 
					"pair", price.Symbol, 
					"price", price.Price)
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func (h *TCPHandler) parsePriceData(data []byte) (*domain.PriceTick, error) {
	var rawData struct {
		Symbol    string  `json:"symbol"`
		Price     float64 `json:"price"`
		Timestamp int64   `json:"timestamp"`
	}
	
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %w", err)
	}
	
	if rawData.Symbol == "" {
		return nil, fmt.Errorf("empty symbol")
	}
	if rawData.Price <= 0 {
		return nil, fmt.Errorf("invalid price: %f", rawData.Price)
	}
	if rawData.Timestamp <= 0 {
		return nil, fmt.Errorf("invalid timestamp: %d", rawData.Timestamp)
	}
	
	timestamp := time.Unix(0, rawData.Timestamp*int64(time.Millisecond))
	
	return &domain.PriceTick{
		Exchange:  h.exchange.Name,
		Symbol:    rawData.Symbol,
		Price:     rawData.Price,
		Ts:        timestamp,
	}, nil
}