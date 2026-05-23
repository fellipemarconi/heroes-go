package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Message string `json:"message"`
}

func Send(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *AppError
	if errors.As(err, &appErr) {
		reportError(w, appErr.Message, "")
		w.WriteHeader(appErr.Status)

		_ = json.NewEncoder(w).Encode(Response{
			Message: appErr.Message,
		})

		return
	}

	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		reportError(w, "validation failed", "")
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

	stack := debug.Stack()
	reportError(w, err.Error(), stackLocation(stack))
	w.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(w).Encode(Response{
		Message: "internal server error",
	})
}

type errorReporter interface {
	SetError(message string, location string)
}

func reportError(w http.ResponseWriter, message string, location string) {
	if reporter, ok := w.(errorReporter); ok {
		reporter.SetError(message, location)
	}
}

func stackLocation(stack []byte) string {
	lines := strings.Split(string(stack), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || !strings.Contains(trimmed, ".go:") {
			continue
		}

		if strings.Contains(trimmed, "/runtime/") ||
			strings.Contains(trimmed, "\\runtime\\") ||
			strings.Contains(trimmed, "/internal/pkg/apierror/") ||
			strings.Contains(trimmed, "\\internal\\pkg\\apierror\\") {
			continue
		}

		fields := strings.Fields(trimmed)
		if len(fields) > 0 {
			return fields[0]
		}
	}

	return ""
}
