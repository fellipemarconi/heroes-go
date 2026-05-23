package hero

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewHeroHandler(t *testing.T) {
	service := &Service{}
	handler := NewHeroHandler(service)

	if handler == nil {
		t.Fatal("NewHeroHandler() returned nil")
	}

	if handler.service != service {
		t.Errorf("NewHeroHandler() service = %p, want %p", handler.service, service)
	}
}

func TestHandler_CreateHero(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/api/hero", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	h.CreateHero(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("CreateHero() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestHandler_GetHero(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/api/hero", nil)
	rec := httptest.NewRecorder()

	h.GetHero(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("GetHero() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandler_ListHeroes(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/api/heroes?is_active=not-bool", nil)
	rec := httptest.NewRecorder()

	h.ListHeroes(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("ListHeroes() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandler_UpdateHeroStatus(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPatch, "/api/hero/status", strings.NewReader(`{"hero_id":"id","is_active":true}`))
	rec := httptest.NewRecorder()

	h.UpdateHeroStatus(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("UpdateHeroStatus() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestHandler_UpdateHeroImage(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/api/hero/image", strings.NewReader("not-multipart"))
	rec := httptest.NewRecorder()

	h.UpdateHeroImage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("UpdateHeroImage() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
