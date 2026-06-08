package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"

	userpb "api_gateway/internal/gen/user/v1"
)

type UserHandler struct {
	client userpb.UserServiceClient
}

func NewUserHandler(client userpb.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) outCtx(r *http.Request) (context.Context, context.CancelFunc) {
	md := metadata.Pairs(
		"x-service-token", os.Getenv("SERVICE_TOKEN"),
		"x-correlation-id", correlationID(r),
	)
	return context.WithTimeout(metadata.NewOutgoingContext(r.Context(), md), 10*time.Second)
}

func (h *UserHandler) pathID(r *http.Request) (int32, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return int32(id), true
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := h.pathID(r)
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.GetProfile(ctx, &userpb.GetProfileRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := h.pathID(r)
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var body struct {
		Username string `json:"username"`
		Phone    string `json:"phone"`
		Fullname string `json:"fullname"`
		Avatar   string `json:"avatar"`
		Gender   int32  `json:"gender"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.UpdateProfile(ctx, &userpb.UpdateProfileRequest{
		Id:       id,
		Username: body.Username,
		Phone:    body.Phone,
		Fullname: body.Fullname,
		Avatar:   body.Avatar,
		Gender:   body.Gender,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := h.pathID(r)
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := h.outCtx(r)
	defer cancel()

	resp, err := h.client.DeactivateUser(ctx, &userpb.DeactivateUserRequest{Id: id})
	if err != nil {
		writeGRPCError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": resp.Id})
}
