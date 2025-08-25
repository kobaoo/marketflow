package exchange

import (
	"context"
	"log/slog"
	"marketflow/internal/domain"
	"math/rand"
	"time"
)

type Message struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

var symbols []string = []string{
	"BTCUSDT",
	"DOGEUSDT",
	"TONUSDT",
	"SOLUSDT",
	"ETHUSDT",
}

var prices map[string]float64 = map[string]float64{
	"BTCUSDT":  117000,
	"DOGEUSDT": 0.232,
	"TONUSDT":  3.47,
	"SOLUSDT":  188,
	"ETHUSDT":  4460,
}

func (r ExchangeClient) startGenerator(ctx context.Context, exchange_name string, out chan<- domain.PriceTick) {
	slog.Info("Starting price generator", "exchange", exchange_name)
	
	ticker := time.NewTicker(100 * time.Millisecond) // Generate prices every 100ms
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down price generator", "exchange", exchange_name)
			return
		case <-ticker.C:
			for _, symbol := range symbols {
				prices[symbol] = r.fluctuateNumber(prices[symbol])
				
				tick := domain.PriceTick{
					Exchange: exchange_name,
					Symbol:   symbol,
					Price:    prices[symbol],
					Ts:       time.Now(),
				}
				
				// Try to send, drop if channel is full
				select {
				case out <- tick:
					// sent successfully
				default:
					slog.Debug("⚠ Dropping stale message", "exchange", exchange_name, "symbol", symbol)
				}
			}
		}
	}
}

func (r ExchangeClient) fluctuateNumber(original float64) float64 {
	rand.NewSource(time.Now().UnixNano())
	fluctuation := (rand.Float64() * 0.4) - 0.2 // generate a random number between -0.2 and 0.2
	return original * (1 + fluctuation)
}
