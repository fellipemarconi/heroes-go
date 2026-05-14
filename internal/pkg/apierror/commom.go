package apierror

import (
	"net/http"
)

type AppError struct {
	Message string
	Status  int
}

func (e *AppError) Error() string {
	return e.Message
}

func New(status int, message string) *AppError {
	return &AppError{
		Status:  status,
		Message: message,
	}
}

var (
	ErrInvalidBody         = New(http.StatusBadRequest, "invalid request body")
	ErrInvalidID           = New(http.StatusBadRequest, "invalid id")
	ErrUserNotFound        = New(http.StatusNotFound, "user not found")
	ErrInvalidCredentials  = New(http.StatusUnauthorized, "invalid credentials")
	ErrEmailAlreadyUsed    = New(http.StatusConflict, "email already used")
	ErrInvalidToken        = New(http.StatusUnauthorized, "invalid token")
	ErrHeroSlugAlreadyUsed = New(http.StatusConflict, "hero slug already used")
)
