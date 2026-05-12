package user

import (
	sqlc "api/internal/infra/db/sqlc"

	"github.com/go-chi/chi/v5"
)

func RoutesAuth(queries *sqlc.Queries, r chi.Router) {
	service := NewUserService(queries)
	handler := NewUserHandler(service)

	r.Post("/api/auth/register", handler.CreateUser)
	r.Post("/api/auth/login", handler.SignInUser)
}

func RoutesUser(queries *sqlc.Queries, r chi.Router) {
	service := NewUserService(queries)
	handler := NewUserHandler(service)

	r.Get("/api/user", handler.GetUser)
	r.Delete("/api/user", handler.DeleteUser)
	r.Patch("/api/user", handler.UpdateUser)
}
