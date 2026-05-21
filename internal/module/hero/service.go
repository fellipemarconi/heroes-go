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

func (s *Service) CreateHero(ctx context.Context, userID string, input *CreateHeroInput) error {
	id, err := types.ParseUUID(userID)
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

func (s *Service) GetHeroBySlug(ctx context.Context, slug string) (*Hero, error) {
	hero, err := s.queries.GetHeroBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.ErrHeroNotFound
		}
		return nil, err
	}

	return &Hero{
		ID:          hero.ID.String(),
		UserID:      hero.UserID.String(),
		Name:        hero.Name,
		Slug:        hero.Slug,
		Alignment:   hero.Alignment,
		Universe:    hero.Universe,
		Powers:      hero.Powers,
		Description: hero.Description.String,
	}, nil
}

func (s *Service) ListHeroes(ctx context.Context, input *ListHeroesInput) ([]*Hero, error) {
	params := sqlc.ListHeroesParams{
		Universe: pgtype.Text{
			String: input.Universe,
			Valid:  input.Universe != "",
		},
		Alignment: pgtype.Text{
			String: input.Alignment,
			Valid:  input.Alignment != "",
		},
		PageOffset: input.PageOffset,
		PageSize:   input.PageSize,
	}
	if input.IsActive != nil {
		params.IsActive = pgtype.Bool{Bool: *input.IsActive, Valid: true}
	}

	heroes, err := s.queries.ListHeroes(ctx, params)
	if err != nil {
		return nil, err
	}

	heroesList := make([]*Hero, len(heroes))
	for i, hero := range heroes {
		heroesList[i] = &Hero{
			ID:          hero.ID.String(),
			UserID:      hero.UserID.String(),
			Name:        hero.Name,
			Slug:        hero.Slug,
			Alignment:   hero.Alignment,
			Universe:    hero.Universe,
			Powers:      hero.Powers,
			Description: hero.Description.String,
		}
	}

	return heroesList, nil
}
