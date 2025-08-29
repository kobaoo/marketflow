package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"marketflow/internal/domain"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domain.Repository {
	return &Repository{db: db}
}

func (r *Repository) StoreMinAgg(ctx context.Context, aggs []*domain.MinuteAgg) error {
	if len(aggs) == 0 {
		return fmt.Errorf("no data to save")
	}

	query := `
		INSERT INTO minute_prices (exchange, pair_name, timestamp, average_price, min_price, max_price)
		VALUES %s
	`

	valueStrings := make([]string, 0, len(aggs))
	valueArgs := make([]interface{}, 0, len(aggs)*6)

	for i, agg := range aggs {
		start := i*6 + 1
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)", start, start+1, start+2, start+3, start+4, start+5),
		)

		valueArgs = append(valueArgs,
			agg.Exchange,
			agg.Symbol,
			agg.Ts,
			agg.Avg,
			agg.Min,
			agg.Max,
		)
	}

	stmt := fmt.Sprintf(query, strings.Join(valueStrings, ","))
	_, err := r.db.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		slog.Error("Error storing data", "err", err)
		return err
	}

	return nil
}

func (r *Repository) Ping(ctx context.Context) error {
	if r == nil || r.db == nil {
		return errors.New("repository/db is nil")
	}

	var one int
	return r.db.QueryRowContext(ctx, "SELECT 1").Scan(&one)
}

func (r *Repository) GetHighestPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	var price float64

	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(max_price)
		FROM minute_prices
		WHERE pair_name = $1
		GROUP BY pair_name
		`, symbol).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetHighestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) (float64, error) {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(max_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2
		GROUP BY pair_name, exchange
		`, symbol, exchange).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetHighestPriceBySymbolAndPeriod(
	ctx context.Context, symbol string, period time.Duration,
) (float64, error) {
	cutoff := time.Now().UTC().Add(-period)

	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(max_price)
		FROM minute_prices
		WHERE pair_name = $1
		  AND "timestamp" > $2
		GROUP BY pair_name
	`, symbol, cutoff).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetLowestPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MIN(min_price)
		FROM minute_prices
		WHERE pair_name = $1
		GROUP BY pair_name
		`, symbol).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetHighestPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	cutoff := time.Now().UTC().Add(-period)

	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MAX(max_price)
		FROM minute_prices
		WHERE pair_name = $1
		  AND exchange  = $2
		  AND "timestamp" > $3
		GROUP BY pair_name, exchange
	`, symbol, exchange, cutoff).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetLowestPriceBySymbolAndExchange(
	ctx context.Context, symbol, exchange string,
) (float64, error) {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MIN(min_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2
		GROUP BY pair_name, exchange
	`, symbol, exchange).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetLowestPriceBySymbolAndPeriod(
	ctx context.Context, symbol string, period time.Duration,
) (float64, error) {
	cutoff := time.Now().UTC().Add(-period)

	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MIN(min_price)
		FROM minute_prices
		WHERE pair_name = $1
		  AND "timestamp" > $2
		GROUP BY pair_name
	`, symbol, cutoff).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetLowestPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	cutoff := time.Now().UTC().Add(-period)

	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT MIN(min_price)
		FROM minute_prices
		WHERE pair_name = $1
		  AND exchange  = $2
		  AND "timestamp" > $3
		GROUP BY pair_name, exchange
	`, symbol, exchange, cutoff).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetAvgPriceBySymbol(
	ctx context.Context, symbol string,
) (float64, error) {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT AVG(average_price)
		FROM minute_prices
		WHERE pair_name = $1
		GROUP BY pair_name
	`, symbol).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetAvgPriceBySymbolAndExchange(
	ctx context.Context, symbol, exchange string,
) (float64, error) {
	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT AVG(average_price)
		FROM minute_prices
		WHERE pair_name = $1 AND exchange = $2
		GROUP BY pair_name, exchange
	`, symbol, exchange).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetAvgPriceBySymbolAndPeriod(
	ctx context.Context, symbol string, period time.Duration,
) (float64, error) {
	cutoff := time.Now().UTC().Add(-period)

	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT AVG(average_price)
		FROM minute_prices
		WHERE pair_name = $1
		  AND "timestamp" > $2
		GROUP BY pair_name
	`, symbol, cutoff).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}

func (r *Repository) GetAvgPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	cutoff := time.Now().UTC().Add(-period)

	var price float64
	err := r.db.QueryRowContext(ctx, `
		SELECT AVG(average_price)
		FROM minute_prices
		WHERE pair_name = $1
		  AND exchange  = $2
		  AND "timestamp" > $3
		GROUP BY pair_name, exchange
	`, symbol, exchange, cutoff).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		slog.Error("Repository error", "err", err)
		return 0, err
	}
	return price, nil
}
