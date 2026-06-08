package handler

import (
	"net/http"

	authpb "api_gateway/internal/gen/auth/v1"
	otppb "api_gateway/internal/gen/otp/v1"
	userpb "api_gateway/internal/gen/user/v1"
)

func NewRouter(
	authClient authpb.AuthServiceClient,
	userClient userpb.UserServiceClient,
	otpClient otppb.OtpServiceClient,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health/live", healthOK)
	mux.HandleFunc("GET /health/ready", healthOK)

	auth := NewAuthHandler(authClient)
	mux.HandleFunc("POST /v1/auth/register", auth.Register)
	mux.HandleFunc("POST /v1/auth/login", auth.Login)
	mux.HandleFunc("POST /v1/auth/refresh", auth.Refresh)
	mux.HandleFunc("POST /v1/auth/verify-account", auth.VerifyAccount)
	mux.HandleFunc("POST /v1/auth/password-reset", auth.ResetPassword)

	usr := NewUserHandler(userClient)
	mux.HandleFunc("GET /v1/user/{id}", usr.GetProfile)
	mux.HandleFunc("PATCH /v1/user/{id}", usr.UpdateProfile)
	mux.HandleFunc("DELETE /v1/user/{id}", usr.DeactivateUser)

	otp := NewOtpHandler(otpClient)
	mux.HandleFunc("POST /v1/otp/request", otp.Request)
	mux.HandleFunc("POST /v1/otp/verify", otp.Verify)

	return stripInternalHeaders(mux)
}

func healthOK(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func stripInternalHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range []string{"X-User-ID", "X-Roles", "X-Internal-Service", "X-Internal-User", "X-Forwarded-User"} {
			r.Header.Del(h)
		}
		if r.Header.Get("X-Correlation-ID") == "" {
			r.Header.Set("X-Correlation-ID", newID())
		}
		next.ServeHTTP(w, r)
	})
}
