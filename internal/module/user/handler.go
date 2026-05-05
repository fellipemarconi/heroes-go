package user

import (
	"api/internal/pkg/apierror"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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
