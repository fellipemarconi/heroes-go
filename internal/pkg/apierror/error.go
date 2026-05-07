package apierror

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Errors  map[string]string `json:"errors,omitempty"`
	Message string            `json:"message,omitempty"`
	Code    int               `json:"-"`
}

func (e *APIError) Error() string { return e.Message }

func New(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

func fromValidation(err validator.ValidationErrors) *APIError {
	errs := make(map[string]string)
	for _, e := range err {
		errs[e.Field()] = e.Tag()
	}
	return &APIError{Code: http.StatusUnprocessableEntity, Errors: errs}
}

func Send(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		w.WriteHeader(apiErr.Code)
		json.NewEncoder(w).Encode(apiErr)
		return
	}

	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		Send(w, fromValidation(valErr))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&APIError{Message: "internal server error"})
}
