package user

import (
	"api/internal/pkg/apierror"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewUserHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.Send(w, apierror.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	if err := h.service.CreateUser(r.Context(), &input); err != nil {
		apierror.Send(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
