package domain

import "time"

type (
	Symbol   string // "BTCUSDT",
	Exchange string // "ex1", "ex2", "ex3", "test-1", ...
)

type PriceTick struct {
	Exhange Exchange
	Symbol  Symbol
	Price   float64
	Ts      time.Time
}

type MinuteAgg struct {
	Exchange Exchange
	Symbol   Symbol
	Ts       time.Time // start or end minute window
	Avg      float64
	Min      float64
	Max      float64
}

type HealthStatus struct {
	Status     string            `json:"status"` // ok|degraded|down
	Mode       string            `json:"mode"`   // live|test
	RedisOK    bool              `json:"redis_ok"`
	PostgresOK bool              `json:"postgres_ok"`
	Sources    map[Exchange]bool `json:"sources"` // биржа -> работает ли listener
}
