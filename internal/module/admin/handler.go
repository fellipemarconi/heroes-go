package admin

import (
	"api/internal/pkg/apierror"
	"encoding/json"
	"log"
	"net/http"
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
