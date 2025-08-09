package domain

import (
	"context"
	"time"
)

// интерфейсы (ports): ExchangeClient, Cache, Repository, Aggregator, ModeController

type ExhangeClient interface {
	// Starts ticker reading. Should be able to auto-reconnect inside itself
	Start(ctx context.Context) (<-chan PriceTick, <-chan error)
	Name() Exchange
}

type Cache interface {
	SetLatest(ctx context.Context, ex Exchange, sym Symbol, price float64, ttlSec int) error
	GetLatest(ctx context.Context, ex *Exchange, sym Symbol) (float64, bool, error) // ex=nil -> по всем

	// Window 60с: adding and reading of range; also period clean
	AppendToWindow(ctx context.Context, ex Exchange, sym Symbol, tick PriceTick) error
	ReadWindow(ctx context.Context, ex Exchange, sym Symbol, since time.Time) ([]float64, error)
	TrimOld(ctx context.Context, olderThan time.Time) error

	// for /highest|/lowest|/average?period=...
	ReadWindowPeriod(ctx context.Context, ex *Exchange, sym Symbol, since time.Time) ([]float64, error)
}

type Repository interface {
	InitSchema(ctx context.Context) error
	InsertMinuteAggBatch(ctx context.Context, rows []MinuteAgg) error
	// Fallback for API, if Redis unavailable
	LastAgg(ctx context.Context, ex *Exchange, sym Symbol, limit int) ([]MinuteAgg, error)
}

type ModeController interface {
	SwitchToLive(ctx context.Context) error
	SwitchToTest(ctx context.Context) error
	CurrentMode() string
}

type Aggregator interface {
	// Calculates avg/min/max by array of prices
	Aggregate(prices []float64, ex Exchange, sym Symbol, ts time.Time) MinuteAgg
}

// UseCases for API

type MarketService interface {
	// Latest
	GetLatestPrice(ctx context.Context, sym Symbol, ex *Exchange) (float64, error)
	// Highest/Lowest/Average за период
	GetHighest(ctx context.Context, sym Symbol, ex *Exchange, period time.Duration) (float64, error)
	GetLowest(ctx context.Context, sym Symbol, ex *Exchange, period time.Duration) (float64, error)
	GetAverage(ctx context.Context, sym Symbol, ex *Exchange, period time.Duration) (float64, error)

	// Mode
	SetModeLive(ctx context.Context) error
	SetModeTest(ctx context.Context) error
	GetMode() string

	// Health
	GetHealth(ctx context.Context) (HealthStatus, error)
}
