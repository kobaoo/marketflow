package domain

import "time"

type PriceTick struct {
	Exhange string
	Symbol  string
	Price   float64
	Ts      time.Time
}

type MinuteAgg struct {
	Exchange string
	Symbol   string
	Ts       time.Time // start or end minute window
	Avg      float64
	Min      float64
	Max      float64
}

type HealthStatus struct {
	Status     string          `json:"status"` // ok|degraded|down
	Mode       string          `json:"mode"`   // live|test
	RedisOK    bool            `json:"redis_ok"`
	PostgresOK bool            `json:"postgres_ok"`
	Sources    map[string]bool `json:"sources"` // биржа -> работает ли listener
}
