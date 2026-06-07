package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var (
	ErrConflict           = errors.New("credential already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type Service struct {
	mu      sync.Mutex
	nextID  int
	nowFunc func() time.Time
	users   map[string]RegistrationResponse
	devices map[string]string
}

func NewService(nowFunc func() time.Time) *Service {
	return &Service{nextID: 1, nowFunc: nowFunc, users: make(map[string]RegistrationResponse), devices: make(map[string]string)}
}

func (s *Service) Register(req RegisterRequest) (RegistrationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[req.Email]; ok {
		return RegistrationResponse{}, ErrConflict
	}
	token, err := randomToken()
	if err != nil {
		return RegistrationResponse{}, err
	}
	res := RegistrationResponse{ID: s.nextID, Email: req.Email, Token: token, ExpiresAtToken: s.nowFunc().Add(30 * time.Minute)}
	s.nextID++
	s.users[req.Email] = res
	return res, nil
}

func (s *Service) Login(req LoginRequest) (LoginResponse, error) {
	return LoginResponse{}, ErrInvalidCredentials
}

func (s *Service) Refresh(req RefreshRequest) (LoginResponse, error) {
	return LoginResponse{}, ErrInvalidToken
}

func (s *Service) VerifyAccount(req VerifyAccountRequest) (LoginResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[req.Email]
	if !ok || user.ID != req.UserID || user.Token != req.VerifiedToken || !user.ExpiresAtToken.After(s.nowFunc()) {
		return LoginResponse{}, ErrInvalidToken
	}
	return LoginResponse{}, ErrInvalidToken
}

func (s *Service) ResetPassword(req ResetPasswordRequest) (ResetPasswordResponse, error) {
	return ResetPasswordResponse{}, ErrInvalidToken
}

func randomToken() (string, error) {
	var data [32]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(data[:]), nil
}
