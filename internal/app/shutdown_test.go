package app

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestShutdownEmptyApp(t *testing.T) {
	if err := (&App{}).Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown() error = %v, want nil", err)
	}
}

func TestShutdownReturnsCacheCloseError(t *testing.T) {
	closedCache := redis.NewClient(&redis.Options{})
	_ = closedCache.Close()

	err := (&App{Cache: closedCache}).Shutdown(context.Background())
	if err == nil {
		t.Fatal("Shutdown() error = nil, want close error")
	}
	if !errors.Is(err, redis.ErrClosed) {
		t.Fatalf("Shutdown() error = %v, want redis.ErrClosed", err)
	}
}
