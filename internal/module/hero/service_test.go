package hero

import (
	"api/internal/pkg/apierror"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewHeroService(t *testing.T) {
	service := NewHeroService(nil, nil)
	if service == nil {
		t.Fatal("NewHeroService() returned nil")
	}
	if service.validate == nil {
		t.Fatal("NewHeroService() validate is nil")
	}
}

func TestService_CreateHero(t *testing.T) {
	service := NewHeroService(nil, nil)

	err := service.CreateHero(context.Background(), "invalid-id", &CreateHeroInput{})
	if !errors.Is(err, apierror.ErrInvalidID) {
		t.Fatalf("CreateHero() error = %v, want %v", err, apierror.ErrInvalidID)
	}
}

func TestService_GetHeroBySlug(t *testing.T) {
	t.Skip("requires database")
}

func TestService_ListHeroes(t *testing.T) {
	t.Skip("requires database")
}

func TestService_UpdateHeroStatus(t *testing.T) {
	service := NewHeroService(nil, nil)
	validID := uuid.NewString()

	t.Run("invalid hero id", func(t *testing.T) {
		err := service.UpdateHeroStatus(context.Background(), validID, "invalid-id", true)
		if !errors.Is(err, apierror.ErrInvalidID) {
			t.Fatalf("UpdateHeroStatus() error = %v, want %v", err, apierror.ErrInvalidID)
		}
	})

	t.Run("invalid user id", func(t *testing.T) {
		err := service.UpdateHeroStatus(context.Background(), "invalid-id", validID, true)
		if !errors.Is(err, apierror.ErrInvalidID) {
			t.Fatalf("UpdateHeroStatus() error = %v, want %v", err, apierror.ErrInvalidID)
		}
	})
}

func TestService_UpdateHeroImage(t *testing.T) {
	service := NewHeroService(nil, nil)
	validID := uuid.NewString()

	t.Run("invalid hero id", func(t *testing.T) {
		err := service.UpdateHeroImage(context.Background(), validID, "invalid-id", &UpdateHeroImageInput{})
		if !errors.Is(err, apierror.ErrInvalidID) {
			t.Fatalf("UpdateHeroImage() error = %v, want %v", err, apierror.ErrInvalidID)
		}
	})

	t.Run("invalid user id", func(t *testing.T) {
		err := service.UpdateHeroImage(context.Background(), "invalid-id", validID, &UpdateHeroImageInput{})
		if !errors.Is(err, apierror.ErrInvalidID) {
			t.Fatalf("UpdateHeroImage() error = %v, want %v", err, apierror.ErrInvalidID)
		}
	})
}
