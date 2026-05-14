package user

import (
	"api/internal/pkg/apierror"
	"context"
	"errors"
	"testing"
)

func TestNewUserService(t *testing.T) {
	service := NewUserService(nil, nil)
	if service == nil {
		t.Fatal("NewUserService() returned nil")
	}

	if service.db != nil {
		t.Errorf("NewUserService() db = %v, want nil", service.db)
	}

	if service.queries != nil {
		t.Errorf("NewUserService() queries = %v, want nil", service.queries)
	}

	if service.validate == nil {
		t.Error("NewUserService() validate = nil, want non-nil")
	}
}

func TestService_CreateUser(t *testing.T) {
	s := NewUserService(nil, nil)

	err := s.CreateUser(context.Background(), &CreateUserInput{})
	if err == nil {
		t.Fatal("Service.CreateUser() expected validation error, got nil")
	}
}

func TestService_SignInUser(t *testing.T) {
	s := NewUserService(nil, nil)

	_, err := s.SignInUser(context.Background(), &SignInUserInput{})
	if err == nil {
		t.Fatal("Service.SignInUser() expected validation error, got nil")
	}
}

func TestService_GetUser(t *testing.T) {
	s := NewUserService(nil, nil)

	_, err := s.GetUser(context.Background(), "invalid-id")
	if !errors.Is(err, apierror.ErrInvalidID) {
		t.Errorf("Service.GetUser() error = %v, want %v", err, apierror.ErrInvalidID)
	}
}

func TestService_DeleteUser(t *testing.T) {
	s := NewUserService(nil, nil)

	err := s.DeleteUser(context.Background(), "invalid-id")
	if !errors.Is(err, apierror.ErrInvalidID) {
		t.Errorf("Service.DeleteUser() error = %v, want %v", err, apierror.ErrInvalidID)
	}
}

func TestService_UpdateUser(t *testing.T) {
	s := NewUserService(nil, nil)

	err := s.UpdateUser(context.Background(), "invalid-id", &UpdateUserInput{})
	if !errors.Is(err, apierror.ErrInvalidBody) {
		t.Errorf("Service.UpdateUser() error = %v, want %v", err, apierror.ErrInvalidBody)
	}
}

func TestService_ForgotPassword(t *testing.T) {
	s := NewUserService(nil, nil)

	err := s.ForgotPassword(context.Background(), &ForgotPasswordInput{})
	if err == nil {
		t.Fatal("Service.ForgotPassword() expected validation error, got nil")
	}
}

func TestService_ResetPassword(t *testing.T) {
	s := NewUserService(nil, nil)

	err := s.ResetPassword(context.Background(), &ResetPasswordInput{})
	if err == nil {
		t.Fatal("Service.ResetPassword() expected validation error, got nil")
	}
}

func TestService_UpdateProfileImage(t *testing.T) {
	s := NewUserService(nil, nil)

	err := s.UpdateProfileImage(context.Background(), "invalid-id", &UpdateProfileImageInput{})
	if !errors.Is(err, apierror.ErrInvalidID) {
		t.Errorf("Service.UpdateProfileImage() error = %v, want %v", err, apierror.ErrInvalidID)
	}
}
