package user

import (
	sqlc "api/internal/infra/db/sqlc"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(queries *sqlc.Queries, r chi.Router) {
	service := NewUserService(queries)
	handler := NewUserHandler(service)

	r.Post("/api/auth/register", handler.CreateUser)
	r.Post("/api/auth/login", handler.SignInUser)
}
