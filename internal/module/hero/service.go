package hero

import (
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/infra/storage"
	"api/internal/pkg/apierror"
	"api/internal/pkg/types"
	"context"
	"encoding/json"
	"errors"
	"log"

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
		IsActive:    hero.IsActive,
		Image:       hero.Image.String,
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
			IsActive:    hero.IsActive,
			Image:       hero.Image.String,
		}
	}

	return heroesList, nil
}

func (s *Service) UpdateHeroStatus(ctx context.Context, userID string, heroID string, isActive bool) error {
	hID, err := types.ParseUUID(heroID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	uID, err := types.ParseUUID(userID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	rows, err := s.queries.SetHeroActiveStatus(ctx, sqlc.SetHeroActiveStatusParams{
		ID:       hID,
		UserID:   uID,
		IsActive: isActive,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return apierror.ErrHeroNotFound
	}

	return nil
}

func (s *Service) UpdateHeroImage(
	ctx context.Context,
	userID string,
	heroID string,
	input *UpdateHeroImageInput,
) error {
	hID, err := types.ParseUUID(heroID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	uID, err := types.ParseUUID(userID)
	if err != nil {
		return apierror.ErrInvalidID
	}

	err = storage.ValidateImageFile(
		input.FileHeader,
		input.ContentType,
	)
	if err != nil {
		return err
	}

	path, err := storage.UploadFile(
		ctx,
		input.File,
		input.FileHeader.Size,
		input.ContentType,
		"image",
	)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		_ = storage.DeleteFile(ctx, path)
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := s.queries.WithTx(tx)

	hero, err := qtx.GetHeroByID(ctx, sqlc.GetHeroByIDParams{
		ID:     hID,
		UserID: uID,
	})
	if err != nil {
		_ = storage.DeleteFile(ctx, path)

		if errors.Is(err, pgx.ErrNoRows) {
			return apierror.ErrHeroNotFound
		}

		return err
	}

	oldImage := ""
	if hero.Image.Valid {
		oldImage = hero.Image.String
	}

	err = qtx.UpdateHeroImage(
		ctx,
		sqlc.UpdateHeroImageParams{
			ID:     hID,
			UserID: uID,
			Image: pgtype.Text{
				String: path,
				Valid:  true,
			},
		},
	)
	if err != nil {
		_ = storage.DeleteFile(ctx, path)
		return err
	}

	metadata, err := json.Marshal(map[string]any{
		"size":          input.FileHeader.Size,
		"mime_type":     input.ContentType,
		"original_name": input.FileHeader.Filename,
	})
	if err != nil {
		_ = storage.DeleteFile(ctx, path)
		return err
	}

	err = qtx.CreateFile(
		ctx,
		sqlc.CreateFileParams{
			ID:       pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Path:     path,
			Type:     sqlc.FileTypeImage,
			Metadata: metadata,
			UserID:   uID,
		},
	)
	if err != nil {
		_ = storage.DeleteFile(ctx, path)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		_ = storage.DeleteFile(ctx, path)
		return err
	}

	if oldImage != "" {
		_ = storage.DeleteFile(ctx, oldImage)

		_ = s.queries.DeleteFileByPath(
			ctx,
			sqlc.DeleteFileByPathParams{
				Path:   oldImage,
				UserID: uID,
			},
		)
	}

	return nil
}

func (s *Service) SearchHeroes(ctx context.Context, input *SearchHeroesInput) ([]SearchHeroResponse, int64, error) {
	if err := s.validate.Struct(input); err != nil {
		return nil, 0, err
	}

	universeValue := ""
	if input.Universe != "" {
		universeValue = input.Universe
	}
	alignmentValue := ""
	if input.Alignment != "" {
		alignmentValue = input.Alignment
	}

	heroes, err := s.queries.SearchHeroes(ctx, sqlc.SearchHeroesParams{
		PlaintoTsquery: input.Query,
		Column2:        universeValue,
		Column3:        alignmentValue,
		Limit:          input.PageSize,
		Offset:         input.PageOffset,
	})
	log.Print(heroes)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.queries.CountSearchHeroes(ctx, sqlc.CountSearchHeroesParams{
		PlaintoTsquery: input.Query,
		Column2:        universeValue,
		Column3:        alignmentValue,
	})
	if err != nil {
		return nil, 0, err
	}

	results := make([]SearchHeroResponse, len(heroes))
	for i, hero := range heroes {
		description := ""
		if hero.Description.Valid {
			description = hero.Description.String
		}
		image := ""
		if hero.Image.Valid {
			image = hero.Image.String
		}
		results[i] = SearchHeroResponse{
			ID:          hero.ID.String(),
			Name:        hero.Name,
			Slug:        hero.Slug,
			Universe:    hero.Universe,
			Alignment:   hero.Alignment,
			Powers:      hero.Powers,
			Description: description,
			Image:       image,
			Rank:        hero.Rank,
		}
	}

	return results, total, nil
}
