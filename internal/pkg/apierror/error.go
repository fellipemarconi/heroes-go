package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
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

	w.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(w).Encode(Response{
		Message: "internal server error",
	})
}
