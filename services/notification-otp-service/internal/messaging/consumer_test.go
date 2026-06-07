package messaging

import (
	"errors"
	"testing"
	"time"
)

func TestWithRetryReturnsAfterSuccess(t *testing.T) {
	attempts := 0
	err := withRetry(3, time.Nanosecond, func() error {
		attempts++
		if attempts == 1 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("withRetry() error = %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestWithRetryReturnsLastError(t *testing.T) {
	want := errors.New("failed")
	err := withRetry(2, time.Nanosecond, func() error { return want })
	if !errors.Is(err, want) {
		t.Fatalf("withRetry() error = %v, want %v", err, want)
	}
}
