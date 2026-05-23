package storage

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadFile(
	ctx context.Context,
	file multipart.File,
	size int64,
	contentType string,
	folder string,
) (string, error) {
	bucket := os.Getenv("MINIO_BUCKET")

	filename := uuid.NewString()

	path := fmt.Sprintf(
		"%s/%s",
		folder,
		filename,
	)

	_, err := Client.PutObject(
		ctx,
		bucket,
		path,
		file,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", err
	}

	return path, nil
}

func DeleteFile(
	ctx context.Context,
	path string,
) error {
	bucket := os.Getenv("MINIO_BUCKET")

	return Client.RemoveObject(
		ctx,
		bucket,
		path,
		minio.RemoveObjectOptions{},
	)
}

var (
	ErrInvalidFileSize      = errors.New("invalid file size")
	ErrInvalidFileMimeType  = errors.New("invalid file mime type")
	ErrInvalidFileExtension = errors.New("invalid file extension")
	ErrFileNotFound         = errors.New("file not found")
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxFileSize = 5 << 20 // 5MB

func ValidateImageFile(
	header *multipart.FileHeader,
	contentType string,
) error {
	if header.Size > maxFileSize {
		return ErrInvalidFileSize
	}

	if !allowedMimeTypes[contentType] {
		return ErrInvalidFileMimeType
	}

	ext := strings.ToLower(
		filepath.Ext(header.Filename),
	)

	if !allowedExtensions[ext] {
		return ErrInvalidFileExtension
	}

	return nil
}

func DownloadFile(
	ctx context.Context,
	path string,
) (*minio.Object, minio.ObjectInfo, error) {
	bucket := os.Getenv("MINIO_BUCKET")

	object, err := Client.GetObject(
		ctx,
		bucket,
		path,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}

	info, err := object.Stat()
	if err != nil {
		_ = object.Close()

		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" || errResp.Code == "NoSuchObject" {
			return nil, minio.ObjectInfo{}, ErrFileNotFound
		}

		return nil, minio.ObjectInfo{}, err
	}

	return object, info, nil
}
