package user

import (
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

type ProfileResponse struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	HiddenPhoneNumber string `json:"hidden_phone_number"`
	FullName          string `json:"fullname"`
	HiddenEmail       string `json:"hidden_email"`
	Avatar            string `json:"avatar"`
	Gender            int    `json:"gender"`
	TwoFactorEnabled  bool   `json:"two_factor_enabled"`
	IsActive          bool   `json:"is_active"`
	CreatedAt         string `json:"created_at"`
}

func NewClient(baseURL string, serviceToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), serviceToken: serviceToken, httpClient: httpClient}
}

func (c *Client) GetProfile(ctx context.Context, correlationID string, id int) (ProfileResponse, error) {
	var response ProfileResponse
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/internal/users/%d", c.baseURL, id), nil)
	if err != nil {
		return ProfileResponse{}, fmt.Errorf("create user profile request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.serviceToken)
	req.Header.Set("X-Correlation-ID", correlationID)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return ProfileResponse{}, fmt.Errorf("send user profile request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ProfileResponse{}, fmt.Errorf("user service status %d", res.StatusCode)
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return ProfileResponse{}, fmt.Errorf("decode user profile response: %w", err)
	}
	return response, nil
}
