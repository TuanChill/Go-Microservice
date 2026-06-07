package user

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestMemoryRepositoryProfileLifecycle(t *testing.T) {
	now := time.Date(2026, 6, 7, 10, 0, 0, 0, time.UTC)
	repo := NewMemoryRepository(func() time.Time { return now })

	identity, err := repo.Create("user@example.com")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if identity.ID != 1 || identity.Email != "user@example.com" {
		t.Fatalf("Create() = %#v", identity)
	}

	profile, err := repo.UpdateProfile(identity.ID, UpdateProfileRequest{Username: "tuanchill", Phone: "84901234567", FullName: "Tuan", Avatar: "avatar.png", Gender: 1})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	if profile.Username != "tuanchill" || profile.HiddenPhoneNumber != "84****67" || profile.FullName != "Tuan" || profile.Gender != 1 {
		t.Fatalf("UpdateProfile() = %#v", profile)
	}

	profile, err = repo.GetProfile(identity.ID)
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if !profile.CreatedAt.Equal(now) || !profile.IsActive {
		t.Fatalf("GetProfile() = %#v", profile)
	}

	if _, err := repo.Deactivate(identity.ID); err != nil {
		t.Fatalf("Deactivate() error = %v", err)
	}
	_, err = repo.GetProfile(identity.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetProfile() after deactivate error = %v, want ErrNotFound", err)
	}
}

func TestMemoryRepositoryRejectsDuplicateEmail(t *testing.T) {
	repo := NewMemoryRepository(time.Now)
	if _, err := repo.Create("user@example.com"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	_, err := repo.Create("user@example.com")
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("duplicate Create() error = %v, want ErrConflict", err)
	}
}

func TestMemoryRepositoryConcurrentAccess(t *testing.T) {
	repo := NewMemoryRepository(time.Now)
	identity, err := repo.Create("user@example.com")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = repo.GetProfile(identity.ID)
			_, _ = repo.UpdateProfile(identity.ID, UpdateProfileRequest{Username: "tuanchill", Phone: "84901234567", FullName: "Tuan", Avatar: "avatar.png", Gender: 1})
		}()
	}
	wg.Wait()

	if _, err := repo.Deactivate(identity.ID); err != nil {
		t.Fatalf("Deactivate() error = %v", err)
	}
}
