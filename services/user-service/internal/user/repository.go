package user

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrConflict = errors.New("user already exists")
	ErrNotFound = errors.New("user not found")
)

type Repository interface {
	Create(email string) (Identity, error)
	GetProfile(id int) (Profile, error)
	UpdateProfile(id int, req UpdateProfileRequest) (Profile, error)
	Deactivate(id int) (DestroyAccountResponse, error)
}

type MemoryRepository struct {
	mu      sync.Mutex
	nextID  int
	users   map[int]Profile
	emails  map[string]int
	nowFunc func() time.Time
}

func NewMemoryRepository(nowFunc func() time.Time) *MemoryRepository {
	return &MemoryRepository{nextID: 1, users: make(map[int]Profile), emails: make(map[string]int), nowFunc: nowFunc}
}

func (r *MemoryRepository) Create(email string) (Identity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.emails[email]; ok {
		return Identity{}, ErrConflict
	}
	id := r.nextID
	r.nextID++
	r.users[id] = Profile{ID: id, Email: email, IsActive: true, CreatedAt: r.nowFunc()}
	r.emails[email] = id
	return Identity{ID: id, Email: email}, nil
}

func (r *MemoryRepository) GetProfile(id int) (Profile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	profile, ok := r.users[id]
	if !ok || !profile.IsActive {
		return Profile{}, ErrNotFound
	}
	return profile, nil
}

func (r *MemoryRepository) UpdateProfile(id int, req UpdateProfileRequest) (Profile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	profile, ok := r.users[id]
	if !ok || !profile.IsActive {
		return Profile{}, ErrNotFound
	}
	if req.Username != "" {
		profile.Username = req.Username
	}
	if req.Phone != "" {
		profile.Phone = req.Phone
		profile.HiddenPhoneNumber = hidePhoneNumber(req.Phone)
	}
	if req.FullName != "" {
		profile.FullName = req.FullName
	}
	if req.Avatar != "" {
		profile.Avatar = req.Avatar
	}
	profile.Gender = req.Gender
	r.users[id] = profile
	return profile, nil
}

func (r *MemoryRepository) Deactivate(id int) (DestroyAccountResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	profile, ok := r.users[id]
	if !ok || !profile.IsActive {
		return DestroyAccountResponse{}, ErrNotFound
	}
	profile.IsActive = false
	r.users[id] = profile
	return DestroyAccountResponse{ID: id}, nil
}

func hidePhoneNumber(phone string) string {
	if len(phone) <= 4 {
		return phone
	}
	return phone[:2] + "****" + phone[len(phone)-2:]
}
