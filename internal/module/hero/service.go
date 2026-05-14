package hero

import (
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/pkg/apierror"
	"api/internal/pkg/types"
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db       *pgxpool.Pool
	queries  *sqlc.Queries
	validate *validator.Validate
}

func NewHeroService(db *pgxpool.Pool, queries *sqlc.Queries) *Service {
	return &Service{
		db:       db,
		queries:  queries,
		validate: validator.New(),
	}
}

func (s *Service) CreateHero(ctx context.Context, userId string, input *CreateHeroInput) error {
	id, err := types.ParseUUID(userId)
	if err != nil {
		return apierror.ErrInvalidID
	}

	if err = s.validate.Struct(input); err != nil {
		return err
	}

	_, err = s.queries.GetHeroBySlug(ctx, input.Slug)
	if err == nil {
		return apierror.ErrHeroSlugAlreadyUsed
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if err = s.queries.CreateHero(ctx, sqlc.CreateHeroParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:    id,
		Name:      input.Name,
		Slug:      input.Slug,
		Alignment: input.Alignment,
		Universe:  input.Universe,
		Powers:    input.Powers,
		Description: pgtype.Text{
			String: input.Description,
			Valid:  input.Description != "",
		},
	}); err != nil {
		return err
	}

	return nil
}
