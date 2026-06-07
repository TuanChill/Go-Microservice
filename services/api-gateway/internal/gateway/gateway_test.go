package gateway

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGatewayRoutesByPrefix(t *testing.T) {
	auth := targetServer(t, "auth")
	user := targetServer(t, "user")
	otp := targetServer(t, "otp")
	legacy := targetServer(t, "legacy")
	defer auth.Close()
	defer user.Close()
	defer otp.Close()
	defer legacy.Close()

	handler := New([]Route{
		{Prefix: "/v1/auth", Target: mustParse(t, auth.URL)},
		{Prefix: "/v1/user", Target: mustParse(t, user.URL)},
		{Prefix: "/v1/otp", Target: mustParse(t, otp.URL)},
	}, mustParse(t, legacy.URL))

	for _, tc := range []struct {
		path string
		want string
	}{
		{"/v1/auth/register", "auth"},
		{"/v1/user/profile/42", "user"},
		{"/v1/otp/verify", "otp"},
		{"/v1/authz", "legacy"},
		{"/v1/user-settings", "legacy"},
		{"/v1/other", "legacy"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			handler.ServeHTTP(res, req)
			if res.Body.String() != tc.want {
				t.Fatalf("body = %q, want %q", res.Body.String(), tc.want)
			}
		})
	}
}

func TestGatewayHealthEndpoints(t *testing.T) {
	legacy := targetServer(t, "legacy")
	defer legacy.Close()

	handler := New(nil, mustParse(t, legacy.URL))

	for _, path := range []string{"/health/live", "/health/ready"} {
		t.Run(path, func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			handler.ServeHTTP(res, req)
			if res.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
			}
			if res.Body.String() != "ok" {
				t.Fatalf("body = %q, want ok", res.Body.String())
			}
		})
	}
}

func TestGatewayStripsSpoofedInternalHeaders(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, header := range []string{"X-User-ID", "X-User-Email", "X-Roles", "X-Service-Name", "X-Authenticated-User", "X-Internal-User", "X-Internal-Service", "X-Forwarded-User", "X-Forwarded-Roles"} {
			if r.Header.Get(header) != "" {
				t.Fatalf("%s was forwarded", header)
			}
		}
		if r.Header.Get("X-Correlation-ID") != "corr-12345678" {
			t.Fatalf("X-Correlation-ID = %q", r.Header.Get("X-Correlation-ID"))
		}
		if r.Header.Get("X-Request-ID") == "" {
			t.Fatal("X-Request-ID was not set")
		}
	}))
	defer target.Close()

	handler := New([]Route{{Prefix: "/v1/auth", Target: mustParse(t, target.URL)}}, nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/register", nil)
	for _, header := range []string{"X-User-ID", "X-User-Email", "X-Roles", "X-Service-Name", "X-Authenticated-User", "X-Internal-User", "X-Internal-Service", "X-Forwarded-User", "X-Forwarded-Roles"} {
		req.Header.Set(header, "spoofed")
	}
	req.Header.Set("X-Correlation-ID", "corr-12345678")
	handler.ServeHTTP(res, req)
}

func targetServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
}

func mustParse(t *testing.T, raw string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	return parsed
}
