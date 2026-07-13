package cache

import (
	"context"
	"fmt"

	"github.com/harshbarnawa/mintok/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg config.Config) (*redis.Client, error) {
	options, err := NewRedisOptions(cfg)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}

func NewRedisOptions(cfg config.Config) (*redis.Options, error) {
	options, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	return options, nil
}
