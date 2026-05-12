package storage

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func ConnectMinIO() error {
	_ = godotenv.Load()

	client, err := minio.New(
		os.Getenv("MINIO_ENDPOINT"),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				os.Getenv("MINIO_ROOT_USER"),
				os.Getenv("MINIO_ROOT_PASSWORD"),
				"",
			),
			Secure: false,
		},
	)
	if err != nil {
		return err
	}

	Client = client

	ctx := context.Background()

	bucket := os.Getenv("MINIO_BUCKET")

	exists, err := Client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if !exists {
		err = Client.MakeBucket(
			ctx,
			bucket,
			minio.MakeBucketOptions{},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
