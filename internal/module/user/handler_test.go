package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewUserHandler(t *testing.T) {
	service := &Service{}
	handler := NewUserHandler(service)

	if handler == nil {
		t.Fatal("NewUserHandler() returned nil")
	}

	if handler.service != service {
		t.Errorf("NewUserHandler() service = %p, want %p", handler.service, service)
	}
}

func TestHandler_CreateUser(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.CreateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("CreateUser() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandler_SignInUser(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/users/sign-in", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.SignInUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("SignInUser() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandler_GetUser(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	rec := httptest.NewRecorder()

	h.GetUser(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("GetUser() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodDelete, "/users/me", nil)
	rec := httptest.NewRecorder()

	h.DeleteUser(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("DeleteUser() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPatch, "/users/me", strings.NewReader(`{"name":"john"}`))
	rec := httptest.NewRecorder()

	h.UpdateUser(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("UpdateUser() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestHandler_ForgotPassword(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/users/forgot-password", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.ForgotPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("ForgotPassword() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandler_ResetPassword(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/users/reset-password", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.ResetPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("ResetPassword() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandler_UpdateProfileImage(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/users/profile/image", strings.NewReader("not-multipart"))
	rec := httptest.NewRecorder()

	h.UpdateProfileImage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("UpdateProfileImage() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
