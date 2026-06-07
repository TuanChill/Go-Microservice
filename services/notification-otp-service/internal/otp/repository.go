package otp

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("otp not found")

type Repository interface {
	Create(record Record) (Record, error)
	Claim(code string, now time.Time) (Record, error)
}

type MemoryRepository struct {
	mu      sync.Mutex
	nextID  int
	records map[string]Record
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{nextID: 1, records: make(map[string]Record)}
}

func (r *MemoryRepository) Create(record Record) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	record.ID = r.nextID
	r.nextID++
	record.IsActive = true
	r.records[record.Code] = record
	return record, nil
}

func (r *MemoryRepository) Claim(code string, now time.Time) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.records[code]
	if !ok || !record.IsActive || !record.ExpiresAt.After(now) {
		return Record{}, ErrNotFound
	}
	record.IsActive = false
	r.records[code] = record
	return record, nil
}
