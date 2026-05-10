package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorType string

const (
	TypeBadRequest   ErrorType = "BAD_REQUEST"
	TypeUnauthorized ErrorType = "UNAUTHORIZED"
	TypeForbidden    ErrorType = "FORBIDDEN"
	TypeNotFound     ErrorType = "NOT_FOUND"
	TypeConflict     ErrorType = "CONFLICT"
	TypeInternal     ErrorType = "INTERNAL"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func BadRequest(msg string) *AppError {
	return &AppError{Type: TypeBadRequest, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Type: TypeUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Type: TypeForbidden, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Type: TypeNotFound, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Type: TypeConflict, Message: msg}
}

func Internal(msg string) *AppError {
	return &AppError{
		Type:    TypeInternal,
		Message: msg,
	}
}

func mapStatus(t ErrorType) int {
	switch t {
	case TypeBadRequest:
		return http.StatusBadRequest
	case TypeUnauthorized:
		return http.StatusUnauthorized
	case TypeForbidden:
		return http.StatusForbidden
	case TypeNotFound:
		return http.StatusNotFound
	case TypeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func Send(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *AppError
	if !errors.As(err, &appErr) {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(AppError{
			Message: ErrInternalServer.Error(),
		})
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(mapStatus(appErr.Type))
	err = json.NewEncoder(w).Encode(AppError{
		Message: appErr.Message,
	})
	if err != nil {
		return
	}
}
