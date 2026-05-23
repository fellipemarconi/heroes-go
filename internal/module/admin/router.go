package admin

import (
	sqlc "api/internal/infra/db/sqlc"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RoutesAdmin(db *pgxpool.Pool, queries *sqlc.Queries, r chi.Router) {
	service := NewAdminService(db, queries)
	handler := NewAdminHandler(service)

	publicKey, err := loadAdminPublicKey()
	if err != nil {
		log.Fatalf("admin auth key error: %v", err)
	}

	r.Route("/api/admin", func(r chi.Router) {
		r.Use(requireAdminSignature(publicKey))

		r.Get("/stats", handler.GetStats)
		r.Get("/users", handler.ListUsers)
		r.Delete("/users/{id}", handler.DeleteUser)
		r.Get("/containers", handler.ListContainers)
		r.Post("/containers/restart", handler.RestartContainers)
	})
}
