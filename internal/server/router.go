package server

import (
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/module/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(queries *sqlc.Queries) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	user.NewUserModule(queries, r)

	return r
}
