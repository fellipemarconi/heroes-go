package server

import (
	"api/internal/infra/auth"
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/module/hero"
	"api/internal/module/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, queries *sqlc.Queries) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	r.Group(func(r chi.Router) {
		user.RoutesAuth(db, queries, r)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))

		user.RoutesUser(db, queries, r)
		hero.RoutesHero(db, queries, r)
	})

	return r
}
