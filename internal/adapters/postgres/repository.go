package postgres

import (
	"context"
	"database/sql"
	"marketflow/internal/domain"
	"time"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domain.Repository {
	return &Repository{db: db}
}

func (r *Repository) StoreMinAgg(ctx context.Context, agg *domain.MinuteAgg) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO minute_prices (exchange, pair_name, timestamp, average_price, min_price, max_price)
		VALUES ($1, $2, $3, $4, $5, $6)
		`, agg.Exchange, agg.Symbol, agg.Ts, agg.Avg, agg.Min, agg.Max)
	if err != nil {
		panic(err)
	}
}

func (r *Repository) GetHighestPriceBySymbol(ctx context.Context, symbol string) float64 {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(max_price)
		FROM minute_prices
		WHERE pair_name = $1
		GROUP BY pair_name
		`, symbol).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(max_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2
		GROUP BY pair_name, exchange
		`, symbol, exchange).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	now := time.Now().Unix()
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(max_price)
		FROM minute_prices
		WHERE pair_name = $1 AND timestamp > $2
		GROUP BY pair_name
		`, symbol, now-int64(period.Seconds())).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	now := time.Now().Unix()
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(max_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2 AND timestamp > $3
		GROUP BY pair_name, exchange
		`, symbol, exchange, now-int64(period.Seconds())).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetLowestPriceBySymbol(ctx context.Context, symbol string) float64 {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MIN(min_price)
		FROM minute_prices
		WHERE pair_name = $1
		GROUP BY pair_name
		`, symbol).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetLowestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MIN(min_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2
		GROUP BY pair_name, exchange
		`, symbol, exchange).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	now := time.Now().Unix()
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MIN(min_price)
		FROM minute_prices
		WHERE pair_name = $1 AND timestamp > $2
		GROUP BY pair_name
		`, symbol, now-int64(period.Seconds())).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	now := time.Now().Unix()
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MIN(min_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2 AND timestamp > $3
		GROUP BY pair_name, exchange
		`, symbol, exchange, now-int64(period.Seconds())).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetAvgPriceBySymbol(ctx context.Context, symbol string) float64 {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT AVG(average_price)
		FROM minute_prices
		WHERE pair_name = $1
		GROUP BY pair_name
		`, symbol).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetAvgPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT AVG(average_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2
		GROUP BY pair_name, exchange
		`, symbol, exchange).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	now := time.Now().Unix()
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT AVG(average_price)
		FROM minute_prices
		WHERE pair_name = $1 AND timestamp > $2
		GROUP BY pair_name
		`, symbol, now-int64(period.Seconds())).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}

func (r *Repository) GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	now := time.Now().Unix()
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT AVG(average_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2 AND timestamp > $3
		GROUP BY pair_name, exchange
		`, symbol, exchange, now-int64(period.Seconds())).Scan(&price)
	if err != nil {
		return 0
	}
	return price
}
