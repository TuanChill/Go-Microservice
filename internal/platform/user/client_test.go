package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGetProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/users/42" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer service-token" {
			t.Fatalf("Authorization = %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Correlation-ID") != "corr-12345678" {
			t.Fatalf("X-Correlation-ID = %s", r.Header.Get("X-Correlation-ID"))
		}
		_, _ = w.Write([]byte(`{"id":42,"username":"tuanchill","email":"user@example.com","phone":"","hidden_phone_number":"","fullname":"Tuan","hidden_email":"","avatar":"","gender":1,"two_factor_enabled":true,"is_active":true,"created_at":"2026-06-07T10:00:00Z"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "service-token", server.Client())
	profile, err := client.GetProfile(context.Background(), "corr-12345678", 42)
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if profile.ID != 42 || profile.Email != "user@example.com" || profile.FullName != "Tuan" || !profile.TwoFactorEnabled {
		t.Fatalf("GetProfile() = %#v", profile)
	}
}

func TestClientGetProfileRejectsUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "service-token", server.Client())
	_, err := client.GetProfile(context.Background(), "corr-12345678", 42)
	if err == nil {
		t.Fatal("GetProfile() error = nil, want error")
	}
}
