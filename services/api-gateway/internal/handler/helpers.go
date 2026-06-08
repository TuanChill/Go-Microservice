package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var fallbackCounter uint64

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeGRPCError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.NotFound:
		http.Error(w, "not found", http.StatusNotFound)
	case codes.AlreadyExists:
		http.Error(w, "conflict", http.StatusConflict)
	case codes.Unauthenticated:
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	case codes.InvalidArgument:
		http.Error(w, "bad request", http.StatusBadRequest)
	case codes.PermissionDenied:
		http.Error(w, "forbidden", http.StatusForbidden)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("fallback-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&fallbackCounter, 1))
	}
	return hex.EncodeToString(b[:])
}

func correlationID(r *http.Request) string {
	if id := r.Header.Get("X-Correlation-ID"); id != "" {
		return id
	}
	return newID()
}
