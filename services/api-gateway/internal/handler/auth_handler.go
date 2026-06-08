package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/metadata"

	authpb "api_gateway/internal/gen/auth/v1"
)

type AuthHandler struct {
	client authpb.AuthServiceClient
}

func NewAuthHandler(client authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{client: client}
}

func (h *AuthHandler) outCtx(r *http.Request) (context.Context, context.CancelFunc) {
	md := metadata.Pairs(
		"x-service-token", os.Getenv("SERVICE_TOKEN"),
		"x-correlation-id", correlationID(r),
	)
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(r.Context(), md), 10*time.Second)
	return ctx, cancel
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.Register(ctx, &authpb.RegisterRequest{Email: body.Email})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"id":               resp.Id,
		"email":            resp.Email,
		"token":            resp.Token,
		"expires_at_token": resp.ExpiresAtToken.AsTime(),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
		DeviceID   string `json:"device_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.Login(ctx, &authpb.LoginRequest{
		Identifier: body.Identifier,
		Password:   body.Password,
		DeviceId:   body.DeviceID,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":           resp.Id,
		"device_id":    resp.DeviceId,
		"email":        resp.Email,
		"access_token": resp.AccessToken,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
		DeviceID     string `json:"device_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.RefreshToken(ctx, &authpb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
		DeviceId:     body.DeviceID,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

func (h *AuthHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID        int    `json:"user_id"`
		VerifiedToken string `json:"verified_token"`
		Email         string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.VerifyAccount(ctx, &authpb.VerifyAccountRequest{
		UserId:        int32(body.UserID),
		VerifiedToken: body.VerifiedToken,
		Email:         body.Email,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": resp.Id})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token    string `json:"token"`
		UserID   int    `json:"user_id"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.ResetPassword(ctx, &authpb.ResetPasswordRequest{
		Token:    body.Token,
		UserId:   int32(body.UserID),
		Password: body.Password,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": resp.Id})
}
