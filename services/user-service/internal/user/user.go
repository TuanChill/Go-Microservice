package user

import "time"

type CreateRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type Identity struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Profile struct {
	ID                int       `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	HiddenPhoneNumber string    `json:"hidden_phone_number"`
	FullName          string    `json:"fullname"`
	HiddenEmail       string    `json:"hidden_email"`
	Avatar            string    `json:"avatar"`
	Gender            int       `json:"gender"`
	TwoFactorEnabled  bool      `json:"two_factor_enabled"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	FullName string `json:"fullname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
}

type DestroyAccountResponse struct {
	ID int `json:"id"`
}
