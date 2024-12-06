package db

import (
	"EasyDev-co/pp_file_upload/internal/config"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"time"
)

// S3Client структура для работы с Minio клиентом
type S3Client struct {
	Client *minio.Client
}

// NewS3Client создает и инициализирует новый S3 клиент
func NewS3Client(cfg config.Config) (*S3Client, error) {
	transport := &http.Transport{
		MaxIdleConns:          100,              // Максимальное количество соединений, которые могут простаивать
		MaxIdleConnsPerHost:   100,              // Максимальное количество соединений на хост
		MaxConnsPerHost:       100,              // Максимальное количество соединений на хост одновременно
		IdleConnTimeout:       90 * time.Second, // Таймаут для бездействующего соединения
		TLSHandshakeTimeout:   10 * time.Second, // Таймаут для TLS-рукопожатия
		ExpectContinueTimeout: 1 * time.Second,  // Таймаут ожидания "100-continue" ответа
	}

	client, err := minio.New(cfg.YCS3Endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(cfg.AwsAccessKeyId, cfg.AwsSecretAccessKey, ""),
		Secure:    true,
		Region:    cfg.YCRegion,
		Transport: transport,
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing Minio client: %v", err)
	}

	return &S3Client{
		Client: client,
	}, nil
}
