package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"go_template/internal/models"

	_ "github.com/lib/pq"
)

type Options struct {
	MaxRetries int
	RetryDelay time.Duration
}

func Connect(cfg models.DatabaseConfig, opts Options) (*sql.DB, error) {
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 10
	}
	if opts.RetryDelay == 0 {
		opts.RetryDelay = 5 * time.Second
	}

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", cfg.Username, cfg.Name, cfg.Password, cfg.Host, cfg.Port)

	var lastErr error
	for i := 0; i < opts.MaxRetries; i++ {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			lastErr = err
			time.Sleep(opts.RetryDelay)
			continue
		}

		if err := db.Ping(); err != nil {
			lastErr = err
			db.Close()
			time.Sleep(opts.RetryDelay)
			continue
		}

		return db, nil
	}

	return nil, fmt.Errorf("connect postgres after %d attempts: %w", opts.MaxRetries, lastErr)
}
