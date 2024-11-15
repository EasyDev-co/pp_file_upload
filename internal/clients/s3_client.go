package clients

import (
	"EasyDev-co/pp_file_upload/internal/config"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Client структура для работы с Minio клиентом
type S3Client struct {
	Client *minio.Client
}

// NewS3Client создает и инициализирует новый S3 клиент
func NewS3Client(cfg config.Config) (*S3Client, error) {
	client, err := minio.New(cfg.YCS3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AwsAccessKeyId, cfg.AwsSecretAccessKey, ""),
		Secure: true,
		Region: cfg.YCRegion,
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing Minio client: %v", err)
	}

	return &S3Client{
		Client: client,
	}, nil
}
