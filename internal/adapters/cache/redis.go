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

func NewRedisClient(config *config.Config) (*RedisClient, error) {
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

func (r *RedisClient) processLastMinute(ctx context.Context, pair, exchange string) {
	key := fmt.Sprintf("%s:%s:prices", pair, exchange)
	now := time.Now().Unix()

	// Get last 60s values
	vals, _ := r.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprint(now - 60),
		Max: fmt.Sprint(now),
	}).Result()

	if len(vals) == 0 {
		fmt.Println("No data for", pair)
		return
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
		Symbol:   pair,
		Exchange: exchange,
		Ts:       time.Now().Truncate(time.Minute),
		Avg:      avg,
		Min:      min,
		Max:      max,
	}

	// Store into Postgres (pseudo-code)
	fmt.Printf("Saving summary: %+v\n", summary)
	// db.Exec("INSERT INTO prices_summary ...", ...)
}
