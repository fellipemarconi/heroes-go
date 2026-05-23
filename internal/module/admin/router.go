package admin

import (
	sqlc "api/internal/infra/db/sqlc"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RoutesAdmin(db *pgxpool.Pool, queries *sqlc.Queries, r chi.Router) {
	service := NewAdminService(db, queries)
	handler := NewAdminHandler(service)

	r.Get("/api/admin/stats", handler.GetStats)
}
