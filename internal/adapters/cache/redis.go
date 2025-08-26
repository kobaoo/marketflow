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

func (r *RedisClient) GetLatestPriceByPattern(ctx context.Context, pattern string) float64 {
	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		val, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Count: 1,
			Max:   "+inf",
		}).Result()
		if err != nil {
			slog.Error("failed to query redis", "err", err)
			continue
		}
		if len(val) > 0 {
			price, _ := strconv.ParseFloat(val[0], 64)
			return price
		}
	}

	if err := iter.Err(); err != nil {
		slog.Error("scan error", "err", err)
	}
	slog.Warn("no matching keys found", "pattern", pattern)
	return 0
}


func (r *RedisClient) GetLatestPriceBySymbol(ctx context.Context, symbol string) float64 {
	return r.GetLatestPriceByPattern(ctx, symbol + ":*:prices")
}

func (r *RedisClient) GetLatestPriceBySymbolAndExchange(ctx context.Context, symbol, exchange string) float64 {
	key := fmt.Sprintf("%s:%s:prices", symbol, exchange)
	val, _ := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Count: 1,
		Max:   "+inf",
	}).Result()
	if len(val) == 0 {
		slog.Error("Error geting data from Redis")
		return 0
	}
	price, _ := strconv.ParseFloat(val[0], 64)
	return price
}

func (r *RedisClient) GetHighestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	// Example pattern: BTCUSDT:*:prices
	pattern := fmt.Sprintf("%s:*:prices", symbol)

	now := time.Now().Unix()
	minScore := fmt.Sprint(now - int64(period.Seconds()))
	maxScore := fmt.Sprint(now)

	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	var highest float64

	for iter.Next(ctx) {
		key := iter.Val()

		vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min: minScore,
			Max: maxScore,
		}).Result()
		if err != nil {
			slog.Error("failed to query redis", "key", key, "err", err)
			continue
		}

		for _, v := range vals {
			price, err := strconv.ParseFloat(v, 64)
			if err != nil {
				slog.Error("invalid price format", "val", v, "err", err)
				continue
			}
			if price > highest {
				highest = price
			}
		}
	}

	if err := iter.Err(); err != nil {
		slog.Error("scan error", "err", err)
	}

	if highest == 0 {
		slog.Warn("no data found in Redis", "pattern", pattern)
	}
	return highest
}


func (r *RedisClient) GetHighestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	key := fmt.Sprintf("%s:%s:prices", symbol, exchange)
	now := time.Now().Unix()
	vals, _ := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprint(now - int64(period.Seconds())),
		Max: fmt.Sprint(now),
	}).Result()
	if len(vals) == 0 {
		slog.Error("Error geting data from Redis")
		return 0
	}
	max, _ := strconv.ParseFloat(vals[0], 64)
	for _, v := range vals[1:] {
		price, _ := strconv.ParseFloat(v, 64)
		if price > max {
			max = price
		}
	}
	return max
}

func (r *RedisClient) GetLowestPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	// Example pattern: BTCUSDT:*:prices
	pattern := fmt.Sprintf("%s:*:prices", symbol)

	now := time.Now().Unix()
	minScore := fmt.Sprint(now - int64(period.Seconds()))
	maxScore := fmt.Sprint(now)

	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	var lowest float64
	firstFound := false

	for iter.Next(ctx) {
		key := iter.Val()

		vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min: minScore,
			Max: maxScore,
		}).Result()
		if err != nil {
			slog.Error("failed to query redis", "key", key, "err", err)
			continue
		}

		for _, v := range vals {
			price, err := strconv.ParseFloat(v, 64)
			if err != nil {
				slog.Error("invalid price format", "val", v, "err", err)
				continue
			}

			if !firstFound {
				lowest = price
				firstFound = true
				continue
			}

			if price < lowest {
				lowest = price
			}
		}
	}

	if err := iter.Err(); err != nil {
		slog.Error("scan error", "err", err)
	}

	if !firstFound {
		slog.Warn("no data found in Redis", "pattern", pattern)
		return 0
	}

	return lowest
}


func (r *RedisClient) GetLowestPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	key := fmt.Sprintf("%s:%s:prices", symbol, exchange)
	now := time.Now().Unix()
	vals, _ := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprint(now - int64(period.Seconds())),
		Max: fmt.Sprint(now),
	}).Result()
	if len(vals) == 0 {
		slog.Error("Error geting data from Redis")
		return 0
	}
	min, _ := strconv.ParseFloat(vals[0], 64)
	for _, v := range vals[1:] {
		price, _ := strconv.ParseFloat(v, 64)
		if price < min {
			min = price
		}
	}
	return min
}

func (r *RedisClient) GetAvgPriceBySymbolAndPeriod(ctx context.Context, symbol string, period time.Duration) float64 {
	// Example pattern: BTCUSDT:*:prices
	pattern := fmt.Sprintf("%s:*:prices", symbol)

	now := time.Now().Unix()
	minScore := fmt.Sprint(now - int64(period.Seconds()))
	maxScore := fmt.Sprint(now)

	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	var sum float64
	var count int

	for iter.Next(ctx) {
		key := iter.Val()

		vals, err := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
			Min: minScore,
			Max: maxScore,
		}).Result()
		if err != nil {
			slog.Error("failed to query redis", "key", key, "err", err)
			continue
		}

		for _, v := range vals {
			price, err := strconv.ParseFloat(v, 64)
			if err != nil {
				slog.Error("invalid price format", "val", v, "err", err)
				continue
			}
			sum += price
			count++
		}
	}

	if err := iter.Err(); err != nil {
		slog.Error("scan error", "err", err)
	}

	if count == 0 {
		slog.Warn("no data found in Redis", "pattern", pattern)
		return 0
	}

	return sum / float64(count)
}

func (r *RedisClient) GetAvgPriceBySymbolAndPeriodAndExchange(ctx context.Context, symbol, exchange string, period time.Duration) float64 {
	key := fmt.Sprintf("%s:%s:prices", symbol, exchange)
	now := time.Now().Unix()
	vals, _ := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprint(now - int64(period.Seconds())),
		Max: fmt.Sprint(now),
	}).Result()
	if len(vals) == 0 {
		slog.Error("Error geting data from Redis")
		return 0
	}
	var sum float64
	for _, v := range vals {
		price, _ := strconv.ParseFloat(v, 64)
		sum += price
	}
	avg := sum / float64(len(vals))
	return avg
}
