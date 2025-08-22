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

func (r *RedisClient) StoreTick(ctx context.Context, exchange, pair string, price float64) {
	key := fmt.Sprintf("%s:%s:prices", pair, exchange)
	now := time.Now().Unix()

	// Add tick
	r.rdb.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: price,
	})

	// Remove older than 60s
	r.rdb.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprint(now-60))
}

func (r *RedisClient) ProcessLastMinute(ctx context.Context, pair, exchange string) *domain.MinuteAgg {
	key := fmt.Sprintf("%s:%s:prices", pair, exchange)
	now := time.Now().Unix()

	// Get last 60s values
	vals, _ := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprint(now - 60),
		Max: fmt.Sprint(now),
	}).Result()

	if len(vals) == 0 {
		slog.Error("No data for", "pair", pair)
		return nil
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
	return &summary
}

func (r *RedisClient) GetLatestPriceBySymbol(ctx context.Context, symbol string) float64 {
	key := fmt.Sprintf("%s:prices", symbol)
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
	key := fmt.Sprintf("%s:prices", symbol)
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
	key := fmt.Sprintf("%s:prices", symbol)
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
	key := fmt.Sprintf("%s:prices", symbol)
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
