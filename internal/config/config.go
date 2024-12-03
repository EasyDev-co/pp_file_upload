package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
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
	return config, nil
}
