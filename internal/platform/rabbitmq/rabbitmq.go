package rabbitmq

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Options struct {
	MaxRetries int
	RetryDelay time.Duration
}

func Connect(dsn string, opts Options) (*amqp.Connection, error) {
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 5
	}
	if opts.RetryDelay == 0 {
		opts.RetryDelay = 2 * time.Second
	}

	var lastErr error
	for i := 0; i < opts.MaxRetries; i++ {
		conn, err := amqp.Dial(dsn)
		if err != nil {
			lastErr = err
			time.Sleep(opts.RetryDelay)
			continue
		}
		return conn, nil
	}

	return nil, fmt.Errorf("connect rabbitmq after %d attempts: %w", opts.MaxRetries, lastErr)
}
