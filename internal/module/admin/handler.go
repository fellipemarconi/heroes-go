package admin

import (
	"api/internal/pkg/apierror"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewAdminHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		apierror.Send(w, err)
		return
	}
	if err = json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("failed to encode stats response: %v", err)
		return
	}
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		apierror.Send(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("failed to encode users response: %v", err)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		apierror.Send(w, apierror.ErrInvalidID)
		return
	}

	if err := h.service.DeleteUser(r.Context(), userID); err != nil {
		apierror.Send(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
