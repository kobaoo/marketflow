package exchange

import (
	"context"
	"encoding/json"
	"log/slog"
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

func StartGenerator(ctx context.Context, exchange chan<- []byte) {
	for {
		for _, symbol := range symbols {
			select {
			case <-ctx.Done():
				close(exchange)
				return
			default:
				prices[symbol] = fluctuateNumber(prices[symbol])
				message, err := json.Marshal(Message{Symbol: symbol, Price: prices[symbol], Timestamp: time.Now().Unix()})
				if err != nil {
					slog.Error("Error marshalling message", "error", err)
				}
				exchange <- message
			}
		}
	}
}

func fluctuateNumber(original float64) float64 {
	rand.NewSource(time.Now().UnixNano())
	fluctuation := (rand.Float64() * 0.4) - 0.2 // generate a random number between -0.2 and 0.2
	return original * (1 + fluctuation)
}
