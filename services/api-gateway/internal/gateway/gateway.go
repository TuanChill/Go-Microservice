package gateway

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

type Route struct {
	Owner  string
	Prefix string
	Target *url.URL
}

var fallbackIDCounter uint64

type Gateway struct {
	routes []Route
}

func New(routes []Route) *Gateway {
	return &Gateway{routes: routes}
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health/live" || r.URL.Path == "/health/ready" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return
	}

	target, owner := g.routeFor(r.URL.Path)
	if target == nil {
		http.NotFound(w, r)
		return
	}
	log.Printf("gateway route owner=%s path=%s upstream=%s", owner, r.URL.Path, target.String())

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		prepareHeaders(req)
	}
	proxy.ServeHTTP(w, r)
}

func (g *Gateway) routeFor(path string) (*url.URL, string) {
	for _, route := range g.routes {
		if matchesPrefix(path, route.Prefix) {
			return route.Target, route.Owner
		}
	}
	return nil, ""
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
