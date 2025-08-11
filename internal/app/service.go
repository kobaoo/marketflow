package app

import (
	"context"
	"marketflow/internal/domain"
	"sort"
	"time"
)

type service struct {
	cache domain.Cache
	repo domain.Repository
	mode string
}

func NewService(c domain.Cache, r domain.Repository) domain.MarketService {
	return &service{cache: c, repo: r, mode:"test"}
}

func (s *service) GetMode() string { return s.mode }

func (s *service) SetModeLive(ctx context.Context) error { s.mode = "live"; return nil }
func (s *service) SetModeTest(ctx context.Context) error { s.mode = "test"; return nil }

func (s *service) GetHealth(ctx context.Context) (domain.HealthStatus, error) {
	// In-memory: assume OK
	return domain.HealthStatus{
		Status:     "ok",
		Mode:       s.mode,
		RedisOK:    true,
		PostgresOK: true,
		Sources:    map[domain.Exchange]bool{},
	}, nil
}

func (s *service) GetLatestPrice(ctx context.Context, sym domain.Symbol, ex *domain.Exchange) (float64, error) {
	price, ok, err := s.cache.GetLatest(ctx, ex, sym)
	if err != nil { return 0, err }
	if !ok {
		// fallback to repo agg
		aggs, err := s.repo.LastAgg(ctx, ex, sym, 1)
		if err != nil { return 0, err }
		if len(aggs) == 0 { return 0, domain.ErrNotFound }
		return aggs[0].Avg, nil
	}
	return price, nil
}

func (s *service) GetHighest(ctx context.Context, sym domain.Symbol, ex *domain.Exchange, period time.Duration) (float64, error) {
	if period <= 0 { return 0, domain.ErrInvalidPeriod }
	since := time.Now().Add(-period)
	var prices []float64
	var err error
	if ex != nil {
		prices, err = s.cache.ReadWindowPeriod(ctx, *ex, sym, since)
	} else {
		prices, err = s.cache.ReadWindowPeriod(ctx, "", sym, since)
	}
	if err != nil { return 0, err }
	if len(prices) == 0 {
		aggs, err := s.repo.LastAgg(ctx, ex, sym, 10)
		if err != nil { return 0, err }
		if len(aggs) == 0 { return 0, domain.ErrNotFound }
		mx := aggs[0].Max
		for _, a := range aggs[1:] { if a.Max > mx { mx = a.Max } }
		return mx, nil
	}
	mx := prices[0]
	for _, p := range prices[1:] { if p > mx { mx = p } }
	return mx, nil
}

func (s *service) GetLowest(ctx context.Context, sym domain.Symbol, ex *domain.Exchange, period time.Duration) (float64, error) {
	if period <= 0 { return 0, domain.ErrInvalidPeriod }
	since := time.Now().Add(-period)
	var prices []float64
	var err error
	if ex != nil {
		prices, err = s.cache.ReadWindowPeriod(ctx, *ex, sym, since)
	} else {
		prices, err = s.cache.ReadWindowPeriod(ctx, "", sym, since)
	}
	if err != nil { return 0, err }
	if len(prices) == 0 {
		aggs, err := s.repo.LastAgg(ctx, ex, sym, 10)
		if err != nil { return 0, err }
		if len(aggs) == 0 { return 0, domain.ErrNotFound }
		mn := aggs[0].Min
		for _, a := range aggs[1:] { if a.Min < mn { mn = a.Min } }
		return mn, nil
	}
	mn := prices[0]
	for _, p := range prices[1:] { if p < mn { mn = p } }
	return mn, nil
}

func (s *service) GetAverage(ctx context.Context, sym domain.Symbol, ex *domain.Exchange, period time.Duration) (float64, error) {
	if period <= 0 { return 0, domain.ErrInvalidPeriod }
	since := time.Now().Add(-period)
	var prices []float64
	var err error
	if ex != nil {
		prices, err = s.cache.ReadWindowPeriod(ctx, *ex, sym, since)
	} else {
		prices, err = s.cache.ReadWindowPeriod(ctx, "", sym, since)
	}
	if err != nil { return 0, err }
	if len(prices) == 0 {
		aggs, err := s.repo.LastAgg(ctx, ex, sym, 10)
		if err != nil { return 0, err }
		if len(aggs) == 0 { return 0, domain.ErrNotFound }
		// среднее от последних агрегаций
		var sum float64
		var n int
		for _, a := range aggs { sum += a.Avg; n++ }
		return sum / float64(n), nil
	}
	sort.Float64s(prices)
	var sum float64
	for _, p := range prices { sum += p }
	return sum / float64(len(prices)), nil
}