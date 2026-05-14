package hero

import (
	"api/internal/pkg/apierror"
	"api/internal/pkg/ctx"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHeroHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateHero(w http.ResponseWriter, r *http.Request) {
	var input CreateHeroInput

	userId, err := ctx.GetUserIDCtx(r)
	if err != nil {
		apierror.Send(w, err)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		apierror.Send(w, apierror.ErrInvalidBody)
		return
	}

	if err = h.service.CreateHero(r.Context(), userId, &input); err != nil {
		apierror.Send(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
