package notification

type EmailVerificationRequest struct {
	UserID         int    `json:"user_id" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Token          string `json:"token" binding:"required"`
	ExpiresAtToken string `json:"expires_at_token" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type PasswordResetRequest struct {
	UserID int    `json:"user_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Token  string `json:"token" binding:"required"`
}

type AcceptedResponse struct {
	IdempotencyKey string `json:"idempotency_key"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AcceptEmailVerification(idempotencyKey string, req EmailVerificationRequest) AcceptedResponse {
	return AcceptedResponse{IdempotencyKey: idempotencyKey}
}

func (s *Service) AcceptPasswordReset(idempotencyKey string, req PasswordResetRequest) AcceptedResponse {
	return AcceptedResponse{IdempotencyKey: idempotencyKey}
}
