package file

import (
	"api/internal/infra/storage"
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestNewFileService(t *testing.T) {
	service := NewFileService()
	if service == nil {
		t.Fatal("NewFileService() returned nil")
	}
}

func TestService_GetFile(t *testing.T) {
	client, err := minio.New(
		"localhost:9000",
		&minio.Options{
			Creds:  credentials.NewStaticV4("test", "test", ""),
			Secure: false,
		},
	)
	if err != nil {
		t.Fatalf("minio.New() error = %v", err)
	}

	previous := storage.Client
	storage.Client = client
	t.Cleanup(func() {
		storage.Client = previous
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	service := NewFileService()
	_, _, err = service.GetFile(ctx, "image/not-found")
	if err == nil {
		t.Fatal("GetFile() expected error")
	}
}
