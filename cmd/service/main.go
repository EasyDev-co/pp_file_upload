package main

import (
	"EasyDev-co/pp_file_upload/internal/config"
	"EasyDev-co/pp_file_upload/internal/db"
	"EasyDev-co/pp_file_upload/internal/handlers"
	"EasyDev-co/pp_file_upload/internal/repository/s3"
	"EasyDev-co/pp_file_upload/internal/services/image"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	AppConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}
	s3client, err := db.NewS3Client(AppConfig)
	if err != nil {
		log.Fatalf("Error creating S3 client: %v", err)
		return
	}

	s3repository := s3.NewS3Repository(s3client, AppConfig)
	imageService := image.NewImageService(s3repository, AppConfig)

	r := router.New()
	r.POST("/v1/files/upload/", handlers.NewUploadFileHandler(imageService, AppConfig).ServeFastHTTP)

	server := &fasthttp.Server{
		Handler:            r.Handler,
		MaxRequestBodySize: int(AppConfig.MaxUploadSize),
	}

	log.Infof("Starting server on :%s...", AppConfig.AppPort)
	if err := server.ListenAndServe(":" + AppConfig.AppPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
