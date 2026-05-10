package user

import (
	"api/internal/infra/auth"
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/pkg/apierror"
	"api/internal/pkg/types"
	"context"

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
		return apierror.BadRequest(err.Error())
	}

	_, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return apierror.Conflict(ErrEmailAlreadyUsed.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return apierror.Internal(apierror.ErrInternalServer.Error())
	}

	_, err = s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: string(hash),
	})
	if err != nil {
		return apierror.Internal(apierror.ErrInternalServer.Error())
	}

	return nil
}

func (s *Service) SignInUser(ctx context.Context, input *SignInUserInput) (string, error) {
	if err := s.validate.Struct(input); err != nil {
		return "", apierror.BadRequest(err.Error())
	}

	user, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return "", apierror.NotFound(ErrUserNotFound.Error())
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		return "", apierror.Unauthorized(apierror.ErrInvalidCredentials.Error())
	}

	_, token, err := auth.TokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID.String(),
	})

	if err != nil {
		return "", apierror.Internal(apierror.ErrInternalServer.Error())
	}

	return token, nil
}

func (s *Service) GetUser(ctx context.Context, userId string) (*User, error) {
	id, err := types.ParseUUID(userId)
	if err != nil {
		return nil, apierror.BadRequest(apierror.ErrInvalidID.Error())
	}

	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, apierror.NotFound(ErrUserNotFound.Error())
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
		return apierror.BadRequest(apierror.ErrInvalidID.Error())
	}

	_, err = s.queries.GetUserByID(ctx, id)
	if err != nil {
		return apierror.NotFound(ErrUserNotFound.Error())
	}

	err = s.queries.DeleteUser(ctx, id)
	if err != nil {
		return apierror.Internal(apierror.ErrInternalServer.Error())
	}

	return nil
}
