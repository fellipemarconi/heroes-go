package user

import (
	sqlc "api/internal/infra/db/sqlc"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
	})
}

func NewUserModule(queries *sqlc.Queries, r chi.Router) {
	service := NewUserService(queries)
	handler := NewUserHandler(service)

	RegisterRoutes(r, handler)
}
