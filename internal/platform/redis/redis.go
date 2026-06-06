package redis

import (
	"context"
	"fmt"
	"time"

	"go_template/internal/models"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	MaxRetries int
	RetryDelay time.Duration
}

func Connect(ctx context.Context, cfg models.CacheConfig, opts Options) (*redis.Client, error) {
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 10
	}
	if opts.RetryDelay == 0 {
		opts.RetryDelay = 5 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	var lastErr error
	for i := 0; i < opts.MaxRetries; i++ {
		if err := client.Ping(ctx).Err(); err != nil {
			lastErr = err
			time.Sleep(opts.RetryDelay)
			continue
		}
		return client, nil
	}

	client.Close()
	return nil, fmt.Errorf("connect redis after %d attempts: %w", opts.MaxRetries, lastErr)
}
