package auth

import "time"

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RegistrationResponse struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Token          string    `json:"token"`
	ExpiresAtToken time.Time `json:"expires_at_token"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required,min=6"`
	DeviceID   string `json:"device_id" binding:"required"`
}

type LoginResponse struct {
	ID          int    `json:"id"`
	DeviceID    string `json:"device_id"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	DeviceID     string `json:"device_id" binding:"required"`
}

type VerifyAccountRequest struct {
	UserID        int    `json:"user_id" binding:"required"`
	VerifiedToken string `json:"verified_token" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	UserID   int    `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type ResetPasswordResponse struct {
	ID int `json:"id"`
}
