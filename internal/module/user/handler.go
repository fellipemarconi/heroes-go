package user

import (
	"api/internal/pkg/apierror"
	"api/internal/pkg/ctx"
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
		apierror.Send(w, apierror.ErrInvalidBody)
		return
	}

	if err := h.service.CreateUser(r.Context(), &input); err != nil {
		apierror.Send(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) SignInUser(w http.ResponseWriter, r *http.Request) {
	var input SignInUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.Send(w, apierror.ErrInvalidBody)
		return
	}

	token, err := h.service.SignInUser(r.Context(), &input)
	if err != nil {
		apierror.Send(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	}); err != nil {
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := ctx.GetUserIDCtx(r)
	if err != nil {
		apierror.Send(w, err)
		return
	}

	user, err := h.service.GetUser(r.Context(), userId)
	if err != nil {
		apierror.Send(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := ctx.GetUserIDCtx(r)
	if err != nil {
		apierror.Send(w, err)
		return
	}

	if err = h.service.DeleteUser(r.Context(), userId); err != nil {
		apierror.Send(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
