package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"user_service/internal/user"
)

func TestUserEndpoints(t *testing.T) {
	now := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	router := NewRouter(user.NewMemoryRepository(func() time.Time { return now }), "service-token")

	res := performRequest(router, http.MethodPost, "/internal/users", `{"email":"user@example.com"}`)
	if res.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", res.Code, http.StatusCreated)
	}
	if !strings.Contains(res.Body.String(), `"id":1`) {
		t.Fatalf("create body = %s", res.Body.String())
	}

	res = performRequest(router, http.MethodGet, "/internal/users/1", ``)
	if res.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", res.Code, http.StatusOK)
	}
	if !strings.Contains(res.Body.String(), `"email":"user@example.com"`) {
		t.Fatalf("get body = %s", res.Body.String())
	}

	res = performRequest(router, http.MethodPatch, "/internal/users/1", `{"username":"tuanchill","phone":"84901234567","fullname":"Tuan","avatar":"avatar.png","gender":1}`)
	if res.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want %d", res.Code, http.StatusOK)
	}
	if !strings.Contains(res.Body.String(), `"hidden_phone_number":"84****67"`) {
		t.Fatalf("patch body = %s", res.Body.String())
	}

	res = performRequest(router, http.MethodDelete, "/internal/users/1", ``)
	if res.Code != http.StatusOK {
		t.Fatalf("delete status = %d, want %d", res.Code, http.StatusOK)
	}

	res = performRequest(router, http.MethodGet, "/internal/users/1", ``)
	if res.Code != http.StatusNotFound {
		t.Fatalf("get after delete status = %d, want %d", res.Code, http.StatusNotFound)
	}
}

func TestUserEndpointsRequireHeaders(t *testing.T) {
	router := NewRouter(user.NewMemoryRepository(time.Now), "service-token")

	req := httptest.NewRequest(http.MethodGet, "/internal/users/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("missing auth status = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodGet, "/internal/users/1", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("wrong auth status = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodPost, "/internal/users", bytes.NewBufferString(`{"email":"user@example.com"}`))
	req.Header.Set("Authorization", "Bearer service-token")
	req.Header.Set("X-Correlation-ID", "corr-12345678")
	req.Header.Set("Idempotency-Key", "short")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("short idempotency status = %d, want %d", res.Code, http.StatusBadRequest)
	}
}

func TestCreateUserRejectsDuplicateEmail(t *testing.T) {
	router := NewRouter(user.NewMemoryRepository(time.Now), "service-token")
	res := performRequest(router, http.MethodPost, "/internal/users", `{"email":"user@example.com"}`)
	if res.Code != http.StatusCreated {
		t.Fatalf("first create status = %d, want %d", res.Code, http.StatusCreated)
	}
	res = performRequest(router, http.MethodPost, "/internal/users", `{"email":"user@example.com"}`)
	if res.Code != http.StatusConflict {
		t.Fatalf("duplicate create status = %d, want %d", res.Code, http.StatusConflict)
	}
}

func performRequest(handler http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer service-token")
	req.Header.Set("X-Correlation-ID", "corr-12345678")
	req.Header.Set("Idempotency-Key", "idem-12345678901")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res
}
