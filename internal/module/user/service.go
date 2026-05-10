package user

import (
	"api/internal/infra/auth"
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/pkg/apierror"
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

	user, _ := s.queries.GetUserByEmail(ctx, input.Email)
	if user.ID.Valid {
		return apierror.New(http.StatusConflict, "email is already used")
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
		return "", apierror.New(http.StatusNotFound, "User not found")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		return "", apierror.New(http.StatusUnauthorized, "invalid credentials")
	}

	_, token, err := auth.TokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID.String(),
	})

	if err != nil {
		return "", err
	}

	return token, nil
}
