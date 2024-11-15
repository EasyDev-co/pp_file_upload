package services

import (
	"EasyDev-co/pp_file_upload/internal/clients"
	"EasyDev-co/pp_file_upload/internal/config"
	"EasyDev-co/pp_file_upload/internal/dto"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
)

// S3Service структура, которая реализует интерфейс S3Service
type S3Service struct {
	client *clients.S3Client
	cfg    config.Config
}

// NewS3Service создаем новый S3 сервис
func NewS3Service(client *clients.S3Client, cfg config.Config) *S3Service {
	return &S3Service{client: client, cfg: cfg}
}

// BulkUpload метод для загрузки нескольких файлов в S3
func (s *S3Service) BulkUpload(fileReaders []struct {
	Name    string
	Content io.Reader
}) (*dto.UploadResponse, error) {
	ctx := context.Background()
	var fileURLs []string

	for _, file := range fileReaders {
		_, err := s.client.Client.PutObject(
			ctx,
			s.cfg.YCBucketName,
			file.Name,
			file.Content,
			-1,
			minio.PutObjectOptions{ContentType: "application/octet-stream"},
		)
		if err != nil {
			return nil, fmt.Errorf("error uploading file %s to Minio: %v", file.Name, err)
		}

		fileURL := fmt.Sprintf("https://%s/%s/%s", s.cfg.YCS3Endpoint, s.cfg.YCBucketName, file.Name)
		fileURLs = append(fileURLs, fileURL)

	}

	return &dto.UploadResponse{FileURLs: fileURLs}, nil
}

// BulkDelete метод для удаления нескольких файлов из S3
func (s *S3Service) BulkDelete(files []string) error {
	for _, fileName := range files {
		err := s.client.Client.RemoveObject(context.Background(), s.cfg.YCBucketName, fileName, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("error deleting file %s from Minio: %v", fileName, err)
		}
		fmt.Printf("File %s deleted successfully from bucket %s\n", fileName, s.cfg.YCBucketName)
	}

	return nil
}
