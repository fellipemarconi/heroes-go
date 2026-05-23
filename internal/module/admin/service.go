package admin

import (
	sqlc "api/internal/infra/db/sqlc"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db       *pgxpool.Pool
	queries  *sqlc.Queries
	validate *validator.Validate
}

func NewAdminService(db *pgxpool.Pool, queries *sqlc.Queries) *Service {
	return &Service{
		db:       db,
		queries:  queries,
		validate: validator.New(),
	}
}

func (s *Service) GetStats(ctx context.Context) (*Stats, error) {
	stats, err := s.queries.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return &Stats{
		Users:   stats.Users,
		Heroes:  stats.Heroes,
		Images:  stats.Images,
		Storage: stats.StorageBytes,
	}, nil
}
