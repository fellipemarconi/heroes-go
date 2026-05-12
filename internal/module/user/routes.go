package user

import (
	sqlc "api/internal/infra/db/sqlc"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RoutesAuth(db *pgxpool.Pool, queries *sqlc.Queries, r chi.Router) {
	service := NewUserService(db, queries)
	handler := NewUserHandler(service)

	r.Post("/api/auth/register", handler.CreateUser)
	r.Post("/api/auth/login", handler.SignInUser)
	r.Post("/api/auth/forgot-password", handler.ForgotPassword)
	r.Post("/api/auth/reset-password", handler.ResetPassword)
}

func RoutesUser(db *pgxpool.Pool, queries *sqlc.Queries, r chi.Router) {
	service := NewUserService(db, queries)
	handler := NewUserHandler(service)

	r.Get("/api/user", handler.GetUser)
	r.Delete("/api/user", handler.DeleteUser)
	r.Patch("/api/user", handler.UpdateUser)
	r.Post("/api/user/image", handler.UpdateProfileImage)
}
