package hero

import (
	sqlc "api/internal/infra/db/sqlc"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RoutesHero(db *pgxpool.Pool, queries *sqlc.Queries, r chi.Router) {
	service := NewHeroService(db, queries)
	handler := NewHeroHandler(service)

	r.Post("/api/hero", handler.CreateHero)
	r.Get("/api/hero", handler.GetHero)
	r.Get("/api/heroes/search", handler.SearchHeroes)
	r.Get("/api/heroes", handler.ListHeroes)
	r.Patch("/api/hero/status", handler.UpdateHeroStatus)
	r.Post("/api/hero/image", handler.UpdateHeroImage)
}
