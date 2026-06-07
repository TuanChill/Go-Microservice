package otp

import "time"

type Purpose string

const (
	PurposeLogin         Purpose = "login"
	PurposeEmailUpdate   Purpose = "email_update"
	PurposePasswordReset Purpose = "password_reset"
)

type Request struct {
	UserID  int     `json:"user_id" binding:"required"`
	Email   string  `json:"email" binding:"required,email"`
	Purpose Purpose `json:"purpose" binding:"required,oneof=login email_update password_reset"`
}

type Response struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type VerifyRequest struct {
	Otp string `json:"otp" binding:"required"`
}

type VerifyResponse struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

type Record struct {
	ID        int
	UserID    int
	Email     string
	Code      string
	Purpose   Purpose
	IsActive  bool
	ExpiresAt time.Time
}
