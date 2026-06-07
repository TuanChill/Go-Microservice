package auth

import (
	"errors"
	"testing"
	"time"
)

func TestServiceRegister(t *testing.T) {
	now := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	service := NewService(func() time.Time { return now })

	registered, err := service.Register(RegisterRequest{Email: "user@example.com"})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if registered.ID != 1 || registered.Email != "user@example.com" || len(registered.Token) != 64 || !registered.ExpiresAtToken.Equal(now.Add(30*time.Minute)) {
		t.Fatalf("Register() = %#v", registered)
	}
}

func TestServiceRejectsDuplicateRegister(t *testing.T) {
	service := NewService(time.Now)
	if _, err := service.Register(RegisterRequest{Email: "user@example.com"}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	_, err := service.Register(RegisterRequest{Email: "user@example.com"})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("duplicate Register() error = %v, want ErrConflict", err)
	}
}

func TestServiceFailsClosedForCredentialAndTokenFlows(t *testing.T) {
	service := NewService(time.Now)
	if _, err := service.Register(RegisterRequest{Email: "user@example.com"}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	_, err := service.Login(LoginRequest{Identifier: "user@example.com", Password: "secret1", DeviceID: "device-1"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want ErrInvalidCredentials", err)
	}

	_, err = service.Refresh(RefreshRequest{RefreshToken: "refresh-token", DeviceID: "device-1"})
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Refresh() error = %v, want ErrInvalidToken", err)
	}

	_, err = service.VerifyAccount(VerifyAccountRequest{UserID: 1, VerifiedToken: "token", Email: "user@example.com"})
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("VerifyAccount() error = %v, want ErrInvalidToken", err)
	}

	_, err = service.ResetPassword(ResetPasswordRequest{Token: "reset-token", UserID: 1, Password: "secret2"})
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("ResetPassword() error = %v, want ErrInvalidToken", err)
	}
}
