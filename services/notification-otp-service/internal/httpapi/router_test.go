package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"notification_otp_service/internal/notification"
	"notification_otp_service/internal/otp"
)

func TestRequestAndVerifyOtpEndpoints(t *testing.T) {
	now := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	router := NewRouter(
		otp.NewService(otp.NewMemoryRepository(), func() (string, error) { return "123456", nil }, func() time.Time { return now }, 5*time.Minute),
		notification.NewService(),
		"service-token",
	)

	res := performRequest(router, http.MethodPost, "/internal/otp/request", `{"user_id":42,"email":"user@example.com","purpose":"login"}`)
	if res.Code != http.StatusCreated {
		t.Fatalf("request OTP status = %d, want %d", res.Code, http.StatusCreated)
	}
	if !strings.Contains(res.Body.String(), `"user_id":42`) {
		t.Fatalf("request OTP body = %s", res.Body.String())
	}

	res = performRequest(router, http.MethodPost, "/internal/otp/verify", `{"otp":"123456"}`)
	if res.Code != http.StatusOK {
		t.Fatalf("verify OTP status = %d, want %d", res.Code, http.StatusOK)
	}
	if !strings.Contains(res.Body.String(), `"email":"user@example.com"`) {
		t.Fatalf("verify OTP body = %s", res.Body.String())
	}
}

func TestNotificationEndpointsReturnIdempotencyKey(t *testing.T) {
	router := NewRouter(otp.NewService(otp.NewMemoryRepository(), func() (string, error) { return "123456", nil }, time.Now, 5*time.Minute), notification.NewService(), "service-token")

	res := performRequest(router, http.MethodPost, "/internal/notifications/email-verification", `{"user_id":42,"email":"user@example.com","token":"token","expires_at_token":"2026-06-07T10:30:00Z"}`)
	if res.Code != http.StatusAccepted {
		t.Fatalf("email verification status = %d, want %d", res.Code, http.StatusAccepted)
	}
	if !strings.Contains(res.Body.String(), `"idempotency_key":"idem-12345678901"`) {
		t.Fatalf("email verification body = %s", res.Body.String())
	}

	res = performRequest(router, http.MethodPost, "/internal/notifications/password-reset", `{"user_id":42,"email":"user@example.com","token":"token"}`)
	if res.Code != http.StatusAccepted {
		t.Fatalf("password reset status = %d, want %d", res.Code, http.StatusAccepted)
	}
}

func TestInternalEndpointsRequireServiceHeaders(t *testing.T) {
	router := NewRouter(otp.NewService(otp.NewMemoryRepository(), func() (string, error) { return "123456", nil }, time.Now, 5*time.Minute), notification.NewService(), "service-token")
	req := httptest.NewRequest(http.MethodPost, "/internal/otp/request", bytes.NewBufferString(`{}`))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodPost, "/internal/otp/request", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer ")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("empty bearer status = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodPost, "/internal/otp/request", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer wrong-token")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("wrong bearer status = %d, want %d", res.Code, http.StatusUnauthorized)
	}

	req = httptest.NewRequest(http.MethodPost, "/internal/otp/request", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer service-token")
	req.Header.Set("X-Correlation-ID", "short")
	req.Header.Set("Idempotency-Key", "short")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("short headers status = %d, want %d", res.Code, http.StatusBadRequest)
	}
}

func TestEndpointsValidateContractFields(t *testing.T) {
	router := NewRouter(otp.NewService(otp.NewMemoryRepository(), func() (string, error) { return "123456", nil }, time.Now, 5*time.Minute), notification.NewService(), "service-token")

	res := performRequest(router, http.MethodPost, "/internal/otp/request", `{"user_id":42,"email":"user@example.com","purpose":"invalid"}`)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("invalid purpose status = %d, want %d", res.Code, http.StatusBadRequest)
	}

	res = performRequest(router, http.MethodPost, "/internal/notifications/email-verification", `{"user_id":42,"email":"user@example.com","token":"token","expires_at_token":"not-a-date"}`)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("invalid date-time status = %d, want %d", res.Code, http.StatusBadRequest)
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
