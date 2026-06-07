package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"auth_service/internal/auth"
)

func TestRegisterEndpoint(t *testing.T) {
	router := NewRouter(auth.NewService(func() time.Time { return time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC) }), "service-token")

	res := performRequest(router, "/internal/auth/register", `{"email":"user@example.com"}`)
	if res.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", res.Code, http.StatusCreated)
	}
	if !strings.Contains(res.Body.String(), `"email":"user@example.com"`) {
		t.Fatalf("register body = %s", res.Body.String())
	}
}

func TestCredentialAndTokenEndpointsFailClosed(t *testing.T) {
	router := NewRouter(auth.NewService(time.Now), "service-token")

	for _, tc := range []struct {
		name string
		path string
		body string
	}{
		{"login", "/internal/auth/login", `{"identifier":"user@example.com","password":"secret1","device_id":"device-1"}`},
		{"refresh", "/internal/auth/refresh", `{"refresh_token":"refresh-token","device_id":"device-1"}`},
		{"verify", "/internal/auth/verify-account", `{"user_id":1,"verified_token":"verify-token","email":"user@example.com"}`},
		{"reset", "/internal/auth/password-reset", `{"token":"reset-token","user_id":1,"password":"secret2"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := performRequest(router, tc.path, tc.body)
			if res.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestAuthEndpointsRequireServiceHeaders(t *testing.T) {
	router := NewRouter(auth.NewService(time.Now), "service-token")

	req := httptest.NewRequest(http.MethodPost, "/internal/auth/register", bytes.NewBufferString(`{}`))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("missing auth status = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodPost, "/internal/auth/register", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer wrong-token")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("wrong auth status = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodPost, "/internal/auth/register", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer service-token")
	req.Header.Set("X-Correlation-ID", "short")
	req.Header.Set("Idempotency-Key", "short")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("short headers status = %d, want %d", res.Code, http.StatusBadRequest)
	}
}

func performRequest(handler http.Handler, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer service-token")
	req.Header.Set("X-Correlation-ID", "corr-12345678")
	req.Header.Set("Idempotency-Key", "idem-12345678901")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res
}
