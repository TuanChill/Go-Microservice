package config

import (
	"strings"
	"testing"

	"go_template/internal/models"
)

func TestValidateRequiresAPIPort(t *testing.T) {
	err := Validate(models.Config{}, Options{Runtime: RuntimeAPI})
	if err == nil || !strings.Contains(err.Error(), "server port") {
		t.Fatalf("Validate() error = %v, want server port error", err)
	}
}

func TestValidateRequiresEnabledDependencies(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want string
	}{
		{"postgres enabled", Options{EnablePostgres: true}, "database host"},
		{"migrate runtime", Options{Runtime: RuntimeMigrate}, "database host"},
		{"redis enabled", Options{EnableRedis: true}, "cache host"},
		{"rabbitmq enabled", Options{EnableRabbitMQ: true}, "rabbitmq url"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(models.Config{}, tt.opts)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Validate() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestValidateAllowsDisabledDependencies(t *testing.T) {
	cfg := models.Config{}
	cfg.Server.Port = "8000"

	if err := Validate(cfg, Options{Runtime: RuntimeAPI}); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}
