package admin

import (
	"api/internal/infra/containers"
	sqlc "api/internal/infra/db/sqlc"
	"api/internal/infra/logs"
	"api/internal/infra/storage"
	"api/internal/pkg/apierror"
	"api/internal/pkg/types"
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

func (s *Service) ListUsers(ctx context.Context) ([]UserSummary, error) {
	rows, err := s.queries.ListAdminUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]UserSummary, 0, len(rows))
	for _, row := range rows {
		users = append(users, UserSummary{
			ID:        row.ID.String(),
			Email:     row.Email,
			Name:      row.Name,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return users, nil
}

func (s *Service) DeleteUser(ctx context.Context, userId string) error {
	id, err := types.ParseUUID(userId)
	if err != nil {
		return apierror.ErrInvalidID
	}

	paths, err := s.queries.ListFilePathsByUser(ctx, id)
	if err != nil {
		return err
	}

	rows, err := s.queries.AdminDeleteUser(ctx, id)
	if err != nil {
		return err
	}

	if rows == 0 {
		return apierror.ErrUserNotFound
	}

	for _, path := range paths {
		if err := storage.DeleteFile(ctx, path); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ListContainers(ctx context.Context) ([]ContainerSummary, error) {
	containersList, err := containers.List(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]ContainerSummary, 0, len(containersList))
	for _, container := range containersList {
		items = append(items, ContainerSummary{
			Name:   container.Name,
			Status: container.Status,
		})
	}

	return items, nil
}

func (s *Service) RestartContainers(ctx context.Context) error {
	return containers.RestartAll(ctx)
}

func (s *Service) ListLogs(ctx context.Context, limit int) ([]LogEntry, error) {
	entries := logs.List(limit)
	items := make([]LogEntry, 0, len(entries))
	for _, entry := range entries {
		items = append(items, LogEntry{
			Timestamp:  entry.Timestamp,
			Method:     entry.Method,
			Path:       entry.Path,
			Status:     entry.Status,
			DurationMs: entry.DurationMs,
			Error:      entry.Error,
			Location:   entry.Location,
		})
	}

	return items, nil
}
