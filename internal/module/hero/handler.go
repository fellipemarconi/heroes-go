package hero

import (
	"api/internal/pkg/apierror"
	"api/internal/pkg/ctx"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

const defaultPageSize = 20

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

func (h *Handler) GetHero(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		apierror.Send(w, apierror.ErrMissingSlug)
		return
	}

	hero, err := h.service.GetHeroBySlug(r.Context(), slug)
	if err != nil {
		apierror.Send(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(hero); err != nil {
		return
	}
}

func (h *Handler) ListHeroes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	universe := query.Get("universe")
	alignment := query.Get("alignment")

	var isActive *bool
	if value := query.Get("is_active"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			apierror.Send(w, apierror.ErrInvalidQuery)
			return
		}
		isActive = &parsed
	}

	pageSize := defaultPageSize
	if value := query.Get("page_size"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			apierror.Send(w, apierror.ErrInvalidQuery)
			return
		}
		pageSize = parsed
	}

	pageOffset := 0
	if value := query.Get("page_offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			apierror.Send(w, apierror.ErrInvalidQuery)
			return
		}
		pageOffset = parsed
	}

	heroes, err := h.service.ListHeroes(r.Context(), &ListHeroesInput{
		Universe:   universe,
		Alignment:  alignment,
		IsActive:   isActive,
		PageSize:   int32(pageSize),
		PageOffset: int32(pageOffset),
	})
	if err != nil {
		apierror.Send(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(heroes); err != nil {
		return
	}
}
