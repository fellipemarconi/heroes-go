package user

import (
	"api/internal/infra/auth"
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/pkg/apierror"
	"api/internal/pkg/types"
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries  *sqlc.Queries
	validate *validator.Validate
}

func NewUserService(queries *sqlc.Queries) *Service {
	return &Service{
		queries:  queries,
		validate: validator.New(),
	}
}

func (s *Service) CreateUser(ctx context.Context, input *CreateUserInput) error {
	if err := s.validate.Struct(input); err != nil {
		return err
	}

	_, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return apierror.ErrEmailAlreadyUsed
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: string(hash),
	})

	return err
}

func (s *Service) SignInUser(ctx context.Context, input *SignInUserInput) (string, error) {
	if err := s.validate.Struct(input); err != nil {
		return "", err
	}

	user, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apierror.ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", apierror.ErrInvalidCredentials
		}
		return "", err
	}

	_, token, err := auth.TokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID.String(),
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) GetUser(ctx context.Context, userId string) (*User, error) {
	id, err := types.ParseUUID(userId)
	if err != nil {
		return nil, apierror.ErrInvalidID
	}

	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.ErrUserNotFound
		}
		return nil, err
	}

	return &User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, userId string) error {
	id, err := types.ParseUUID(userId)
	if err != nil {
		return apierror.ErrInvalidID
	}

	rows, err := s.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	if rows == 0 {
		return apierror.ErrUserNotFound
	}

	return nil
}
