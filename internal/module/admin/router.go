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

	r.With(requireAdminSignature(publicKey)).Get("/api/admin/stats", handler.GetStats)
}
