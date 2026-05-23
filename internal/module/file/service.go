package file

import (
	"api/internal/infra/storage"
	"context"

	"github.com/minio/minio-go/v7"
)

type Service struct{}

func NewFileService() *Service {
	return &Service{}
}

func (s *Service) GetFile(ctx context.Context, path string) (*minio.Object, minio.ObjectInfo, error) {
	return storage.DownloadFile(ctx, path)
}
