package otp

import (
	"crypto/rand"
	"fmt"
	"time"
)

type CodeGenerator func() (string, error)

type Service struct {
	repo     Repository
	generate CodeGenerator
	now      func() time.Time
	ttl      time.Duration
}

func NewService(repo Repository, generate CodeGenerator, now func() time.Time, ttl time.Duration) *Service {
	return &Service{repo: repo, generate: generate, now: now, ttl: ttl}
}

func (s *Service) Request(req Request) (Response, error) {
	code, err := s.generate()
	if err != nil {
		return Response{}, fmt.Errorf("generate otp: %w", err)
	}

	expiresAt := s.now().Add(s.ttl)
	record, err := s.repo.Create(Record{
		UserID:    req.UserID,
		Email:     req.Email,
		Code:      code,
		Purpose:   req.Purpose,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return Response{}, fmt.Errorf("create otp: %w", err)
	}

	return Response{ID: record.ID, UserID: record.UserID, ExpiresAt: record.ExpiresAt}, nil
}

func (s *Service) Verify(req VerifyRequest) (VerifyResponse, error) {
	record, err := s.repo.Claim(req.Otp, s.now())
	if err != nil {
		return VerifyResponse{}, err
	}
	return VerifyResponse{UserID: record.UserID, Email: record.Email}, nil
}

func GenerateSixDigitCode() (string, error) {
	var bytes [6]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	code := make([]byte, len(bytes))
	for i, b := range bytes {
		code[i] = '0' + b%10
	}
	return string(code), nil
}
