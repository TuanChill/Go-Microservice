package otp

import (
	"errors"
	"testing"
	"time"
)

func TestServiceRequestAndVerify(t *testing.T) {
	now := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	service := NewService(NewMemoryRepository(), func() (string, error) { return "123456", nil }, func() time.Time { return now }, 5*time.Minute)

	created, err := service.Request(Request{UserID: 42, Email: "user@example.com", Purpose: PurposeLogin})
	if err != nil {
		t.Fatalf("Request() error = %v", err)
	}
	if created.ID != 1 || created.UserID != 42 || !created.ExpiresAt.Equal(now.Add(5*time.Minute)) {
		t.Fatalf("Request() = %#v", created)
	}

	verified, err := service.Verify(VerifyRequest{Otp: "123456"})
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if verified.UserID != 42 || verified.Email != "user@example.com" {
		t.Fatalf("Verify() = %#v", verified)
	}

	_, err = service.Verify(VerifyRequest{Otp: "123456"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("second Verify() error = %v, want ErrNotFound", err)
	}
}

func TestServiceVerifyRejectsExpiredOtp(t *testing.T) {
	now := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	repo := NewMemoryRepository()
	_, err := repo.Create(Record{UserID: 42, Email: "user@example.com", Code: "123456", Purpose: PurposeLogin, ExpiresAt: now.Add(-time.Minute)})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	service := NewService(repo, func() (string, error) { return "000000", nil }, func() time.Time { return now }, 5*time.Minute)
	_, err = service.Verify(VerifyRequest{Otp: "123456"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Verify() error = %v, want ErrNotFound", err)
	}
}
