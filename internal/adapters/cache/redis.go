package cache

import (
	"context"
	"fmt"
	"log/slog"
	"marketflow/internal/config"
	"marketflow/internal/domain"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(config *config.Config) (domain.RedisClient, error) {
	ctx := context.Background()

	// Create client
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,     // Redis server address
		Password: config.Redis.Password, // no password set
		DB:       config.Redis.DB,       // use default DB
	})

	// Test connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	slog.Info("Connected to Redis", "pong", pong)
	return &RedisClient{rdb: rdb}, nil
}

func (r *RedisClient) StoreTick(ctx context.Context, exchange, pair string, price float64) error {
	key := fmt.Sprintf("%s:%s:prices", pair, exchange)
	now := time.Now().Unix()

	if err := r.rdb.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: price,
	}).Err(); err != nil {
		return fmt.Errorf("redis ZADD %q: %w", key, err)
	}

	if err := r.rdb.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprint(now-60)).Err(); err != nil {
		return fmt.Errorf("redis ZREMRANGEBYSCORE %q: %w", key, err)
	}

	return nil
}

func (r *RedisClient) ProcessLastMinute(ctx context.Context) []*domain.MinuteAgg {
	exchanges := []string{"ex1", "ex2", "ex3"}
	pairs := []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	var summaries []*domain.MinuteAgg
	for _, exchange := range exchanges {
		for _, pair := range pairs {
			key := fmt.Sprintf("%s:%s:prices", pair, exchange)
			now := time.Now().Unix()

			// Get last 60s values
			vals, _ := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
				Min: fmt.Sprint(now - 60),
				Max: fmt.Sprint(now),
			}).Result()

			if len(vals) == 0 {
				slog.Error("No data for", "exchange", exchange, "pair", pair)
				continue
			}

			// Compute avg, min, max
			var sum, min, max float64
			for i, v := range vals {
				price, _ := strconv.ParseFloat(v, 64)
				sum += price
				if i == 0 || price < min {
					min = price
				}
				if i == 0 || price > max {
					max = price
				}
			}

			avg := sum / float64(len(vals))
			summary := domain.MinuteAgg{
				Exchange: exchange,
				Symbol:   pair,
				Ts:       time.Now().Truncate(time.Minute),
				Avg:      avg,
				Min:      min,
				Max:      max,
			}

			slog.Debug("Counting average for minute", "exchange", exchange, "pair", pair, "avg", avg, "min", min, "max", max)
			summaries = append(summaries, &summary)
		}
	}

	return summaries
}

func (r *RedisClient) GetLatestPriceBySymbol(ctx context.Context, symbol string) (float64, error) {
	pattern := fmt.Sprintf("%s:*:prices", symbol)
	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	var (
		found  bool
		latest int64
		result float64
	)

	for iter.Next(ctx) {
		key := iter.Val()

		vals, err := r.rdb.ZRevRangeWithScores(ctx, key, 0, 0).Result()
		if err != nil {
			return 0, fmt.Errorf("redis query %q: %w", key, err)
		}
		if len(vals) == 0 {
			continue
		}

		score := int64(vals[0].Score)
		price, err := strconv.ParseFloat(vals[0].Member.(string), 64)
		if err != nil {
			return 0, fmt.Errorf("parse price for %q: %w", key, err)
		}

		if !found || score > latest {
			found = true
			latest = score
			result = price
		}
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("scan error: %w", err)
	}

	if !found {
		return 0, domain.ErrNotFound
	}
	return result, nil
}

func (r *RedisClient) GetLatestPriceBySymbolAndExchange(
	ctx context.Context, symbol, exchange string,
) (float64, error) {
	key := fmt.Sprintf("%s:%s:prices", symbol, exchange)
	vals, err := r.rdb.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Max:    "+inf",
		Min:    "-inf",
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("redis ZREVRANGEBYSCORE %q: %w", key, err)
	}
	if len(vals) == 0 {
		return 0, domain.ErrNotFound
	}

	price, err := strconv.ParseFloat(vals[0], 64)
	if err != nil {
		return 0, fmt.Errorf("parse price for %q: %w", key, err)
	}
	return price, nil
}

func (r *RedisClient) GetHighestPriceBySymbolAndPeriod(
	ctx context.Context, symbol string, period time.Duration,
) (float64, error) {
	pattern := fmt.Sprintf("%s:*:prices", symbol)

	now := time.Now().Unix()
	minScore := strconv.FormatInt(now-int64(period.Seconds()), 10)
	maxScore := strconv.FormatInt(now, 10)

	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	const pageSize = 512
	var (
		found    bool
		maxPrice float64
	)

	for iter.Next(ctx) {
		key := iter.Val()
		var offset int64
		for {
			vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
				Min:    minScore,
				Max:    maxScore,
				Offset: offset,
				Count:  pageSize,
			}).Result()
			if err != nil {
				return 0, fmt.Errorf("redis ZRANGEBYSCORE %q: %w", key, err)
			}
			if len(vals) == 0 {
				break
			}

			for _, s := range vals {
				p, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return 0, fmt.Errorf("parse float for key %q: %w", key, err)
				}
				if !found || p > maxPrice {
					found = true
					maxPrice = p
				}
			}

			if len(vals) < pageSize {
				break
			}
			offset += int64(len(vals))
		}
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("scan error for pattern %q: %w", pattern, err)
	}

	if !found {
		return 0, domain.ErrNotFound
	}
	return maxPrice, nil
}

// ===== HIGHEST (symbol + period + exchange) =====

func (r *RedisClient) GetHighestPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	key := fmt.Sprintf("%s:%s:prices", symbol, exchange)

	now := time.Now().Unix()
	minScore := strconv.FormatInt(now-int64(period.Seconds()), 10)
	maxScore := strconv.FormatInt(now, 10)

	const pageSize = 512
	var (
		found  bool
		maxVal float64
		offset int64
	)

	for {
		vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    minScore,
			Max:    maxScore,
			Offset: offset,
			Count:  pageSize,
		}).Result()
		if err != nil {
			return 0, fmt.Errorf("redis ZRANGEBYSCORE %q: %w", key, err)
		}
		if len(vals) == 0 {
			break
		}

		for _, s := range vals {
			p, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, fmt.Errorf("parse float for key %q: %w", key, err)
			}
			if !found || p > maxVal {
				found = true
				maxVal = p
			}
		}

		if len(vals) < pageSize {
			break
		}
		offset += int64(len(vals))
	}

	if !found {
		return 0, domain.ErrNotFound
	}
	return maxVal, nil
}

// ===== LOWEST =====

func (r *RedisClient) GetLowestPriceBySymbolAndPeriod(
	ctx context.Context, symbol string, period time.Duration,
) (float64, error) {
	pattern := fmt.Sprintf("%s:*:prices", symbol)

	now := time.Now().Unix()
	minScore := strconv.FormatInt(now-int64(period.Seconds()), 10)
	maxScore := strconv.FormatInt(now, 10)

	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	const pageSize = 512
	var (
		found  bool
		minVal float64
	)

	for iter.Next(ctx) {
		key := iter.Val()
		var offset int64
		for {
			vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
				Min:    minScore,
				Max:    maxScore,
				Offset: offset,
				Count:  pageSize,
			}).Result()
			if err != nil {
				return 0, fmt.Errorf("redis ZRANGEBYSCORE %q: %w", key, err)
			}
			if len(vals) == 0 {
				break
			}

			for _, s := range vals {
				p, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return 0, fmt.Errorf("parse float for key %q: %w", key, err)
				}
				if !found || p < minVal {
					found = true
					minVal = p
				}
			}

			if len(vals) < pageSize {
				break
			}
			offset += int64(len(vals))
		}
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("scan error for pattern %q: %w", pattern, err)
	}

	if !found {
		return 0, domain.ErrNotFound
	}
	return minVal, nil
}

func (r *RedisClient) GetLowestPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	key := fmt.Sprintf("%s:%s:prices", symbol, exchange)

	now := time.Now().Unix()
	minScore := strconv.FormatInt(now-int64(period.Seconds()), 10)
	maxScore := strconv.FormatInt(now, 10)

	const pageSize = 512
	var (
		found  bool
		minVal float64
		offset int64
	)

	for {
		vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    minScore,
			Max:    maxScore,
			Offset: offset,
			Count:  pageSize,
		}).Result()
		if err != nil {
			return 0, fmt.Errorf("redis ZRANGEBYSCORE %q: %w", key, err)
		}
		if len(vals) == 0 {
			break
		}

		for _, s := range vals {
			p, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, fmt.Errorf("parse float for key %q: %w", key, err)
			}
			if !found || p < minVal {
				found = true
				minVal = p
			}
		}

		if len(vals) < pageSize {
			break
		}
		offset += int64(len(vals))
	}

	if !found {
		return 0, domain.ErrNotFound
	}
	return minVal, nil
}

// ===== AVERAGE =====

func (r *RedisClient) GetAvgPriceBySymbolAndPeriod(
	ctx context.Context, symbol string, period time.Duration,
) (float64, error) {
	pattern := fmt.Sprintf("%s:*:prices", symbol)

	now := time.Now().Unix()
	minScore := strconv.FormatInt(now-int64(period.Seconds()), 10)
	maxScore := strconv.FormatInt(now, 10)

	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	const pageSize = 512
	var (
		sum   float64
		count int64
	)

	for iter.Next(ctx) {
		key := iter.Val()
		var offset int64
		for {
			vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
				Min:    minScore,
				Max:    maxScore,
				Offset: offset,
				Count:  pageSize,
			}).Result()
			if err != nil {
				return 0, fmt.Errorf("redis ZRANGEBYSCORE %q: %w", key, err)
			}
			if len(vals) == 0 {
				break
			}

			for _, s := range vals {
				p, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return 0, fmt.Errorf("parse float for key %q: %w", key, err)
				}
				sum += p
				count++
			}

			if len(vals) < pageSize {
				break
			}
			offset += int64(len(vals))
		}
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("scan error for pattern %q: %w", pattern, err)
	}

	if count == 0 {
		return 0, domain.ErrNotFound
	}
	return sum / float64(count), nil
}

func (r *RedisClient) GetAvgPriceBySymbolAndPeriodAndExchange(
	ctx context.Context, symbol, exchange string, period time.Duration,
) (float64, error) {
	key := fmt.Sprintf("%s:%s:prices", symbol, exchange)

	now := time.Now().Unix()
	minScore := strconv.FormatInt(now-int64(period.Seconds()), 10)
	maxScore := strconv.FormatInt(now, 10)

	const pageSize = 512
	var (
		sum    float64
		count  int64
		offset int64
	)

	for {
		vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    minScore,
			Max:    maxScore,
			Offset: offset,
			Count:  pageSize,
		}).Result()
		if err != nil {
			return 0, fmt.Errorf("redis ZRANGEBYSCORE %q: %w", key, err)
		}
		if len(vals) == 0 {
			break
		}

		for _, s := range vals {
			p, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, fmt.Errorf("parse float for key %q: %w", key, err)
			}
			sum += p
			count++
		}

		if len(vals) < pageSize {
			break
		}
		offset += int64(len(vals))
	}

	if count == 0 {
		return 0, domain.ErrNotFound
	}
	return sum / float64(count), nil
}
