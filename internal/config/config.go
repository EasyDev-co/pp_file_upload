package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	AppPort            string
	AuthSecretKey      string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
	YCBucketName       string
	YCRegion           string
	YCS3Endpoint       string
	WatermarkPath      string
	WatermarkBytes     []byte
	MaxUploadSize      int64
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		err = fmt.Errorf("error loading .env file %v", err)
	}

	config := Config{
		AppPort:            os.Getenv("APP_PORT"),
		AuthSecretKey:      os.Getenv("AUTH_SECRET_KEY"),
		AwsAccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		YCBucketName:       os.Getenv("YC_BUCKET_NAME"),
		YCRegion:           os.Getenv("YC_REGION"),
		YCS3Endpoint:       os.Getenv("YC_S3_ENDPOINT"),
		WatermarkPath:      os.Getenv("WATERMARK_PATH"),
	}
	maxUploadSizeStr := os.Getenv("MAX_UPLOAD_SIZE")
	if maxUploadSizeStr != "" {
		maxUploadSize, err := strconv.ParseInt(maxUploadSizeStr, 10, 64)
		if err != nil {
			return config, fmt.Errorf("error parsing MAX_UPLOAD_SIZE: %v", err)
		}
		config.MaxUploadSize = maxUploadSize
	} else {
		config.MaxUploadSize = 2.5 * 1024 * 1024 * 1024
	}

	if config.WatermarkPath != "" {
		watermarkBytes, err := os.ReadFile(config.WatermarkPath)
		if err != nil {
			return config, fmt.Errorf("error reading watermark file at %s: %v", config.WatermarkPath, err)
		}
		config.WatermarkBytes = watermarkBytes
	}

	return config, nil
}
