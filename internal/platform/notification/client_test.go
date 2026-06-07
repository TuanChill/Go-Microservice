package notification

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientRequestOtp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/otp/request" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer service-token" {
			t.Fatalf("Authorization = %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Correlation-ID") != "corr-12345678" {
			t.Fatalf("X-Correlation-ID = %s", r.Header.Get("X-Correlation-ID"))
		}
		if r.Header.Get("Idempotency-Key") != "idem-12345678901" {
			t.Fatalf("Idempotency-Key = %s", r.Header.Get("Idempotency-Key"))
		}

		var req OtpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		if req.UserID != 42 || req.Email != "user@example.com" || req.Purpose != "login" {
			t.Fatalf("request = %#v", req)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":7,"user_id":42,"expires_at":"2026-06-07T10:05:00Z"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "service-token", server.Client())
	res, err := client.RequestOtp(context.Background(), "corr-12345678", "idem-12345678901", OtpRequest{UserID: 42, Email: "user@example.com", Purpose: "login"})
	if err != nil {
		t.Fatalf("RequestOtp() error = %v", err)
	}
	if res.ID != 7 || res.UserID != 42 || res.ExpiresAt != "2026-06-07T10:05:00Z" {
		t.Fatalf("RequestOtp() = %#v", res)
	}
}

func TestClientRequestOtpRejectsUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(server.URL, "service-token", server.Client())
	_, err := client.RequestOtp(context.Background(), "corr-12345678", "idem-12345678901", OtpRequest{UserID: 42, Email: "user@example.com", Purpose: "login"})
	if err == nil {
		t.Fatal("RequestOtp() error = nil, want error")
	}
}
