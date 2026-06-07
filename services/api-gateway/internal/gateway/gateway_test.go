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
	defer auth.Close()
	defer user.Close()
	defer otp.Close()

	handler := New([]Route{
		{Owner: "auth-service", Prefix: "/v1/auth", Target: mustParse(t, auth.URL)},
		{Owner: "user-service", Prefix: "/v1/user", Target: mustParse(t, user.URL)},
		{Owner: "notification-otp-service", Prefix: "/v1/otp", Target: mustParse(t, otp.URL)},
	})

	for _, tc := range []struct {
		path string
		want string
	}{
		{"/v1/auth/register", "auth"},
		{"/v1/user/profile/42", "user"},
		{"/v1/otp/verify", "otp"},
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

func TestGatewayReportsRouteOwners(t *testing.T) {
	authURL := mustParse(t, "http://auth.local")
	userURL := mustParse(t, "http://user.local")
	handler := New([]Route{
		{Owner: "auth-service", Prefix: "/v1/auth", Target: authURL},
		{Owner: "user-service", Prefix: "/v1/user", Target: userURL},
	})

	for _, tc := range []struct {
		path       string
		wantOwner  string
		wantTarget *url.URL
	}{
		{"/v1/auth/register", "auth-service", authURL},
		{"/v1/user/profile/42", "user-service", userURL},
	} {
		t.Run(tc.path, func(t *testing.T) {
			gotTarget, gotOwner := handler.routeFor(tc.path)
			if gotOwner != tc.wantOwner {
				t.Fatalf("owner = %q, want %q", gotOwner, tc.wantOwner)
			}
			if gotTarget.String() != tc.wantTarget.String() {
				t.Fatalf("target = %q, want %q", gotTarget.String(), tc.wantTarget.String())
			}
		})
	}
}

func TestGatewayReturnsNotFoundForUnownedRoutes(t *testing.T) {
	handler := New([]Route{{Owner: "auth-service", Prefix: "/v1/auth", Target: mustParse(t, "http://auth.local")}})

	for _, path := range []string{"/v1/authz", "/v1/user-settings", "/v1/other"} {
		t.Run(path, func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			handler.ServeHTTP(res, req)
			if res.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusNotFound)
			}
		})
	}
}

func TestGatewayHealthEndpoints(t *testing.T) {
	handler := New(nil)

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

	handler := New([]Route{{Owner: "auth-service", Prefix: "/v1/auth", Target: mustParse(t, target.URL)}})
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
