package domain

import "time"

type PriceTick struct {
	Exchange string
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

type SystemHealth struct {
	Status     string          `json:"status"`      // healthy | degraded | down
	Service    string          `json:"service"`     // marketflow
	Message    string          `json:"message"`     // описание состояния
	Redis      bool            `json:"redis_ok"`    // состояние Redis
	PostgreSQL bool            `json:"postgres_ok"` // состояние Postgres
	Exchanges  map[string]bool `json:"exchanges"`   // имя биржи -> работает ли listener
	Mode       string          `json:"mode"`        // live | test
	Timestamp  time.Time       `json:"timestamp"`   // время проверки
}