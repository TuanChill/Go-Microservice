package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/metadata"

	otppb "api_gateway/internal/gen/otp/v1"
)

type OtpHandler struct {
	client otppb.OtpServiceClient
}

func NewOtpHandler(client otppb.OtpServiceClient) *OtpHandler {
	return &OtpHandler{client: client}
}

func (h *OtpHandler) outCtx(r *http.Request) (context.Context, context.CancelFunc) {
	md := metadata.Pairs(
		"x-service-token", os.Getenv("SERVICE_TOKEN"),
		"x-correlation-id", correlationID(r),
	)
	return context.WithTimeout(metadata.NewOutgoingContext(r.Context(), md), 10*time.Second)
}

func (h *OtpHandler) Request(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID  int    `json:"user_id"`
		Email   string `json:"email"`
		Purpose string `json:"purpose"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.RequestOtp(ctx, &otppb.RequestOtpRequest{
		UserId:  int32(body.UserID),
		Email:   body.Email,
		Purpose: body.Purpose,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         resp.Id,
		"user_id":    resp.UserId,
		"expires_at": resp.ExpiresAt.AsTime(),
	})
}

func (h *OtpHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Otp string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.VerifyOtp(ctx, &otppb.VerifyOtpRequest{Otp: body.Otp})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user_id": resp.UserId,
		"email":   resp.Email,
	})
}
