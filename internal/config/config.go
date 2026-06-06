package config

import (
	"fmt"
	"os"

	"go_template/configs/common/constants"
	"go_template/internal/models"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Runtime string

const (
	RuntimeAPI     Runtime = "api"
	RuntimeWorker  Runtime = "worker"
	RuntimeMigrate Runtime = "migrate"
)

type Options struct {
	Path           string
	Runtime        Runtime
	EnablePostgres bool
	EnableRedis    bool
	EnableRabbitMQ bool
}

func Load(opts Options) (models.Config, error) {
	if opts.Path == "" {
		opts.Path = "configs/yaml"
	}

	_ = godotenv.Load()

	v := viper.New()
	v.AddConfigPath(opts.Path)
	if os.Getenv("ENV") == constants.DevEnvironment {
		v.SetConfigName("config.dev")
	} else {
		v.SetConfigName("config.prod")
	}
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return models.Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg models.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return models.Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := Validate(cfg, opts); err != nil {
		return models.Config{}, err
	}

	return cfg, nil
}

func Validate(cfg models.Config, opts Options) error {
	if opts.Runtime == RuntimeAPI && cfg.Server.Port == "" {
		return fmt.Errorf("server port is required for api runtime")
	}
	if (opts.Runtime == RuntimeMigrate || opts.EnablePostgres) && cfg.Database.Host == "" {
		return fmt.Errorf("database host is required when postgres is enabled")
	}
	if opts.EnableRedis && cfg.Cache.Host == "" {
		return fmt.Errorf("cache host is required when redis is enabled")
	}
	if opts.EnableRabbitMQ && cfg.RabbitMQ.URL == "" {
		return fmt.Errorf("rabbitmq url is required when rabbitmq is enabled")
	}
	return nil
}
