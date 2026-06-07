package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAuthPayloadJSONShape(t *testing.T) {
	payload := Payload{ID: 42, Email: "user@example.com"}

	got := marshalJSON(t, payload)
	want := `{"id":42,"email":"user@example.com"}`
	if got != want {
		t.Fatalf("Payload JSON = %s, want %s", got, want)
	}
}

func TestLoginResponseJSONShape(t *testing.T) {
	response := LoginResponse{
		ID:          42,
		DeviceID:    "device-1",
		Email:       "user@example.com",
		AccessToken: "access-token",
	}

	got := marshalJSON(t, response)
	for _, field := range []string{`"id":42`, `"device_id":"device-1"`, `"email":"user@example.com"`, `"accessToken":"access-token"`} {
		if !strings.Contains(got, field) {
			t.Fatalf("LoginResponse JSON = %s, missing %s", got, field)
		}
	}
}

func TestProfileResponseJSONShape(t *testing.T) {
	response := ProfileResponseJSON{
		ID:               42,
		Email:            "user@example.com",
		TwoFactorEnabled: true,
		IsActive:         true,
	}

	got := marshalJSON(t, response)
	for _, field := range []string{`"id":42`, `"email":"user@example.com"`, `"two_factor_enabled":true`, `"is_active":true`} {
		if !strings.Contains(got, field) {
			t.Fatalf("ProfileResponseJSON JSON = %s, missing %s", got, field)
		}
	}
}

func TestOTPAndVerificationJSONShape(t *testing.T) {
	otpRequest := marshalJSON(t, OtpRequest{Otp: "123456"})
	if otpRequest != `{"otp":"123456"}` {
		t.Fatalf("OtpRequest JSON = %s", otpRequest)
	}

	verification := marshalJSON(t, VerificationResponse{ID: 1, UserId: 42, Token: "token"})
	for _, field := range []string{`"id":1`, `"user_id":42`, `"token":"token"`} {
		if !strings.Contains(verification, field) {
			t.Fatalf("VerificationResponse JSON = %s, missing %s", verification, field)
		}
	}
}

func marshalJSON(t *testing.T, value any) string {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	return string(data)
}
