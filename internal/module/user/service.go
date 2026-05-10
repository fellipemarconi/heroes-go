package user

import (
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

	existing, _ := s.queries.GetUserByEmail(ctx, input.Email)
	if existing.ID.Valid {
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
