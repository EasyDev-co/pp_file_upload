package s3

import (
	"EasyDev-co/pp_file_upload/internal/config"
	"EasyDev-co/pp_file_upload/internal/db"
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

// RepositoryS3 структура, которая реализует интерфейс S3Service
type RepositoryS3 struct {
	client *db.S3Client
	cfg    config.Config
}

// NewS3Repository создаем новый S3 сервис
func NewS3Repository(client *db.S3Client, cfg config.Config) *RepositoryS3 {
	return &RepositoryS3{client: client, cfg: cfg}
}

func (s *RepositoryS3) BulkUpload(fileReaders []struct {
	Name    string
	Content io.Reader
}) (*[]dto.UploadedFilesDTO, error) {
	ctx := context.Background()
	fileURLs := &[]dto.UploadedFilesDTO{}

	type uploadResult struct {
		FileDTO dto.UploadedFilesDTO
		Error   error
	}

	results := make(chan uploadResult, len(fileReaders))
	defer close(results)

	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)

	for _, file := range fileReaders {
		sem <- struct{}{}
		go func(file struct {
			Name    string
			Content io.Reader
		}) {
			defer func() { <-sem }()

			var fileDTO dto.UploadedFilesDTO
			_, err := s.client.Client.PutObject(
				ctx,
				s.cfg.YCBucketName,
				file.Name,
				file.Content,
				-1,
				minio.PutObjectOptions{ContentType: "application/octet-stream"},
			)
			if err != nil {
				results <- uploadResult{Error: fmt.Errorf("error uploading file %s to Minio: %v", file.Name, err)}
				return
			}

			fileURL := fmt.Sprintf("https://%s/%s/%s", s.cfg.YCS3Endpoint, s.cfg.YCBucketName, file.Name)
			fileDTO.FileURL = fileURL
			results <- uploadResult{FileDTO: fileDTO}
		}(file)
	}

	for i := 0; i < len(fileReaders); i++ {
		result := <-results
		if result.Error != nil {
			return nil, result.Error
		}
		*fileURLs = append(*fileURLs, result.FileDTO)
	}

	return fileURLs, nil
}

func (s *RepositoryS3) UploadFile(name string, content io.Reader) (string, error) {
	ctx := context.Background()

	fmt.Printf("URL: %s\n", name)

	fileURL := fmt.Sprintf("https://%s/%s/%s", s.cfg.YCS3Endpoint, s.cfg.YCBucketName, name)
	_, err := s.client.Client.PutObject(
		ctx,
		s.cfg.YCBucketName,
		name,
		content,
		-1,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		log.WithError(err).Error("failed to upload file")
		return "", fmt.Errorf("error uploading file %s: %v", name, err)
	}

	return fileURL, nil
}

// BulkDelete метод для удаления нескольких файлов из S3
func (s *RepositoryS3) BulkDelete(files []string) error {
	for _, fileName := range files {
		err := s.client.Client.RemoveObject(context.Background(), s.cfg.YCBucketName, fileName, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("error deleting file %s from Minio: %v", fileName, err)
		}
		fmt.Printf("File %s deleted successfully from bucket %s\n", fileName, s.cfg.YCBucketName)
	}

	return nil
}
