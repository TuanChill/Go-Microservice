package health

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		path      string
		readiness func() error
		wantCode  int
		wantBody  string
	}{
		{"live", "/health/live", nil, http.StatusOK, `{"status":"ok"}`},
		{"ready nil callback", "/health/ready", nil, http.StatusOK, `{"status":"ready"}`},
		{"ready success", "/health/ready", func() error { return nil }, http.StatusOK, `{"status":"ready"}`},
		{"ready failure", "/health/ready", func() error { return errors.New("down") }, http.StatusServiceUnavailable, `{"status":"unavailable"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			RegisterRoutes(router, tt.readiness)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			router.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantCode)
			}
			if recorder.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", recorder.Body.String(), tt.wantBody)
			}
		})
	}
}
