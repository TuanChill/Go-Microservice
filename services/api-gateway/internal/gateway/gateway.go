package gateway

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

type Route struct {
	Prefix string
	Target *url.URL
}

var fallbackIDCounter uint64

type Gateway struct {
	routes   []Route
	fallback *url.URL
}

func New(routes []Route, fallback *url.URL) *Gateway {
	return &Gateway{routes: routes, fallback: fallback}
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health/live" || r.URL.Path == "/health/ready" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return
	}

	target := g.fallback
	for _, route := range g.routes {
		if matchesPrefix(r.URL.Path, route.Prefix) {
			target = route.Target
			break
		}
	}
	if target == nil {
		http.NotFound(w, r)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		prepareHeaders(req)
	}
	proxy.ServeHTTP(w, r)
}

func prepareHeaders(req *http.Request) {
	for _, header := range []string{"X-User-ID", "X-User-Email", "X-Roles", "X-Service-Name", "X-Authenticated-User", "X-Internal-User", "X-Internal-Service", "X-Forwarded-User", "X-Forwarded-Roles"} {
		req.Header.Del(header)
	}
	if req.Header.Get("X-Correlation-ID") == "" {
		req.Header.Set("X-Correlation-ID", newID())
	}
	if req.Header.Get("X-Request-ID") == "" {
		req.Header.Set("X-Request-ID", newID())
	}
}

func matchesPrefix(path string, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func ParseURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse target URL: %w", err)
	}
	return parsed, nil
}

func newID() string {
	var data [16]byte
	if _, err := rand.Read(data[:]); err != nil {
		return fmt.Sprintf("fallback-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&fallbackIDCounter, 1))
	}
	return hex.EncodeToString(data[:])
}
