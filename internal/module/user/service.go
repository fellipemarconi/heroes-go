package user

import (
	"api/internal/infra/auth"
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/infra/email"
	redisconn "api/internal/infra/redis"
	"api/internal/infra/storage"
	"api/internal/pkg/apierror"
	"api/internal/pkg/types"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db       *pgxpool.Pool
	queries  *sqlc.Queries
	validate *validator.Validate
}

func NewUserService(db *pgxpool.Pool, queries *sqlc.Queries) *Service {
	return &Service{
		db:       db,
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
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
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
		Image:     user.Image.String,
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

func (s *Service) UpdateUser(ctx context.Context, userId string, input *UpdateUserInput) error {
	if input.Name == "" && input.OldPassword == "" && input.NewPassword == "" {
		return apierror.ErrInvalidBody
	}

	if err := s.validate.Struct(input); err != nil {
		return err
	}

	id, err := types.ParseUUID(userId)
	if err != nil {
		return apierror.ErrInvalidID
	}

	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apierror.ErrUserNotFound
		}
		return err
	}

	passwordHash := user.PasswordHash
	name := user.Name

	if input.Name != "" {
		name = input.Name
	}

	if input.OldPassword != "" {
		err = bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(input.OldPassword),
		)
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return apierror.ErrInvalidCredentials
			}
			return err
		}

		var hash []byte

		hash, err = bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		passwordHash = string(hash)
	}

	if err = s.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:           id,
		Name:         name,
		PasswordHash: passwordHash,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Service) ForgotPassword(ctx context.Context, input *ForgotPasswordInput) error {
	if err := s.validate.Struct(input); err != nil {
		return err
	}

	user, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil
	}

	token := uuid.NewString()
	hash := sha256.Sum256([]byte(token))
	key := hex.EncodeToString(hash[:])

	err = redisconn.Client.Set(ctx, "reset:password:"+key, user.ID.String(), 15*time.Minute).Err()
	if err != nil {
		return err
	}

	body, err := email.RenderEmailTemplate(
		"reset-pass.html",
		struct {
			Token string
		}{
			Token: token,
		},
	)
	if err != nil {
		return err
	}

	err = email.SendEmail(
		user.Email,
		"Reset password",
		body,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, input *ResetPasswordInput) error {
	if err := s.validate.Struct(input); err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(input.Token))
	key := hex.EncodeToString(hash[:])

	userId, err := redisconn.Client.Get(ctx, "reset:password:"+key).Result()
	if err != nil {
		return apierror.ErrInvalidToken
	}

	id, err := types.ParseUUID(userId)
	if err != nil {
		return apierror.ErrInvalidID
	}

	hashPass, err := bcrypt.GenerateFromPassword(
		[]byte(input.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	if err = s.queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: string(hashPass),
	}); err != nil {
		return err
	}

	_ = redisconn.Client.Del(ctx, "reset:password:"+key)

	return nil
}

func (s *Service) UpdateProfileImage(
	ctx context.Context,
	userId string,
	input *UpdateProfileImageInput,
) error {
	id, err := types.ParseUUID(userId)
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
		"profile",
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

	user, err := qtx.GetUserByID(ctx, id)
	if err != nil {
		_ = storage.DeleteFile(ctx, path)

		if errors.Is(err, pgx.ErrNoRows) {
			return apierror.ErrUserNotFound
		}

		return err
	}

	oldImage := ""

	if user.Image.Valid {
		oldImage = user.Image.String
	}

	err = qtx.UpdateUserImage(
		ctx,
		sqlc.UpdateUserImageParams{
			ID: id,
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
			Type:     sqlc.FileTypeProfile,
			Metadata: metadata,
			UserID:   id,
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
				UserID: id,
			},
		)
	}

	return nil
}
