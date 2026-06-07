package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseURL      string
	serviceToken string
	httpClient   *http.Client
}

type OtpRequest struct {
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
}

type OtpResponse struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ExpiresAt string `json:"expires_at"`
}

func NewClient(baseURL string, serviceToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), serviceToken: serviceToken, httpClient: httpClient}
}

func (c *Client) RequestOtp(ctx context.Context, correlationID string, idempotencyKey string, req OtpRequest) (OtpResponse, error) {
	var response OtpResponse
	if err := c.post(ctx, "/internal/otp/request", correlationID, idempotencyKey, req, http.StatusCreated, &response); err != nil {
		return OtpResponse{}, err
	}
	return response, nil
}

func (c *Client) post(ctx context.Context, path string, correlationID string, idempotencyKey string, payload any, wantStatus int, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal notification request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create notification request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.serviceToken)
	req.Header.Set("X-Correlation-ID", correlationID)
	req.Header.Set("Idempotency-Key", idempotencyKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send notification request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != wantStatus {
		return fmt.Errorf("notification service status %d", res.StatusCode)
	}
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return fmt.Errorf("decode notification response: %w", err)
	}
	return nil
}
