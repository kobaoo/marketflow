package domain

import "time"

type PriceTick struct {
	Exchange string
	Symbol   string
	Price    float64
	Ts       time.Time
}

type MinuteAgg struct {
	Exchange string
	Symbol   string
	Ts       time.Time
	Avg      float64
	Min      float64
	Max      float64
}

type SystemHealth struct {
	Status     string          `json:"status"`
	Service    string          `json:"service"`
	Message    string          `json:"message"`
	Redis      bool            `json:"redis_ok"`
	PostgreSQL bool            `json:"postgres_ok"`
	Exchanges  map[string]bool `json:"exchanges"`
	Mode       string          `json:"mode"`
	Timestamp  time.Time       `json:"timestamp"`
}

type Key struct {
	Exchange string
	Symbol   string
}
