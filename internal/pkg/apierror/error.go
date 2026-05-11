package apierror

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Message string `json:"message"`
}

func Send(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *AppError
	if errors.As(err, &appErr) {
		w.WriteHeader(appErr.Status)

		_ = json.NewEncoder(w).Encode(Response{
			Message: appErr.Message,
		})

		return
	}

	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		fields := make(map[string]string)

		for _, fieldErr := range validationErr {
			fields[fieldErr.Field()] = fieldErr.Tag()
		}

		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "validation failed",
			"errors":  fields,
		})

		return
	}

	w.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(w).Encode(Response{
		Message: "internal server error",
	})
}
