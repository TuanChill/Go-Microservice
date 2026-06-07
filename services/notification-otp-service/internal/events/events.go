package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrDuplicate = errors.New("duplicate event")

type Envelope struct {
	EventID        string          `json:"event_id"`
	EventType      string          `json:"event_type"`
	EventVersion   int             `json:"event_version"`
	CorrelationID  string          `json:"correlation_id"`
	IdempotencyKey string          `json:"idempotency_key"`
	Producer       string          `json:"producer"`
	OccurredAt     time.Time       `json:"occurred_at"`
	Data           json.RawMessage `json:"data"`
}

type Store interface {
	MarkProcessed(eventID string) error
}

type MemoryStore struct {
	mu        sync.Mutex
	processed map[string]struct{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{processed: make(map[string]struct{})}
}

func (s *MemoryStore) MarkProcessed(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.processed[eventID]; ok {
		return ErrDuplicate
	}
	s.processed[eventID] = struct{}{}
	return nil
}

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Handle(body []byte) error {
	var envelope Envelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("decode event envelope: %w", err)
	}
	if envelope.EventID == "" {
		return errors.New("event_id is required")
	}
	if err := h.store.MarkProcessed(envelope.EventID); err != nil {
		return err
	}
	return nil
}
