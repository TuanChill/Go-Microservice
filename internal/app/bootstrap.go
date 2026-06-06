package app

import (
	"context"
	"database/sql"
	"log/slog"

	"go_template/internal/config"
	"go_template/internal/models"
	"go_template/internal/platform/logger"
	"go_template/internal/platform/postgres"
	"go_template/internal/platform/rabbitmq"
	platformredis "go_template/internal/platform/redis"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config       models.Config
	Logger       *slog.Logger
	DB           *sql.DB
	Cache        *redis.Client
	MessageQueue *amqp.Connection
}

func Bootstrap(ctx context.Context, opts config.Options) (*App, error) {
	if opts.Runtime == config.RuntimeMigrate {
		opts.EnablePostgres = true
	}

	cfg, err := config.Load(opts)
	if err != nil {
		return nil, err
	}

	app := &App{
		Config: cfg,
		Logger: logger.New(),
	}

	if opts.EnablePostgres {
		db, err := postgres.Connect(cfg.Database, postgres.Options{})
		if err != nil {
			return nil, err
		}
		app.DB = db
	}

	if opts.EnableRedis {
		cache, err := platformredis.Connect(ctx, cfg.Cache, platformredis.Options{})
		if err != nil {
			_ = app.Shutdown(ctx)
			return nil, err
		}
		app.Cache = cache
	}

	if opts.EnableRabbitMQ {
		mq, err := rabbitmq.Connect(cfg.RabbitMQ.URL, rabbitmq.Options{})
		if err != nil {
			_ = app.Shutdown(ctx)
			return nil, err
		}
		app.MessageQueue = mq
	}

	return app, nil
}
