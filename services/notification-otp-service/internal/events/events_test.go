package events

import (
	"errors"
	"testing"
)

func TestHandlerRejectsDuplicateEvent(t *testing.T) {
	handler := NewHandler(NewMemoryStore())
	body := []byte(`{
		"event_id":"evt-1",
		"event_type":"otp.requested",
		"event_version":1,
		"correlation_id":"corr-1",
		"idempotency_key":"idem-1",
		"producer":"notification-otp-service",
		"occurred_at":"2026-06-07T10:00:00Z",
		"data":{}
	}`)

	if err := handler.Handle(body); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if err := handler.Handle(body); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate Handle() error = %v, want ErrDuplicate", err)
	}
}

func TestHandlerRequiresEventID(t *testing.T) {
	handler := NewHandler(NewMemoryStore())
	err := handler.Handle([]byte(`{"event_type":"otp.requested"}`))
	if err == nil {
		t.Fatal("Handle() error = nil, want error")
	}
}
