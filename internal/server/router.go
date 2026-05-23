package server

import (
	"api/internal/infra/auth"
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/module/admin"
	"api/internal/module/file"
	"api/internal/module/hero"
	"api/internal/module/user"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(db *pgxpool.Pool, queries *sqlc.Queries) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	r.Group(func(r chi.Router) {
		user.RoutesAuth(db, queries, r)
		admin.RoutesAdmin(db, queries, r)

		setupDocsRoutes(r)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))

		user.RoutesUser(db, queries, r)
		hero.RoutesHero(db, queries, r)
		file.RoutesFile(r)
	})

	return r
}

func setupDocsRoutes(r chi.Router) {
	r.Get("/docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "./docs/openapi.yaml")
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/openapi.yaml"),
	))
}
