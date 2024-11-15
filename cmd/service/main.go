package main

import (
	"EasyDev-co/pp_file_upload/internal/clients"
	"EasyDev-co/pp_file_upload/internal/config"
	"EasyDev-co/pp_file_upload/internal/handlers"
	"EasyDev-co/pp_file_upload/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
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

	s3client, err := clients.NewS3Client(AppConfig)
	if err != nil {
		log.Fatalf("Error creating S3 client: %v", err)
		return
	}

	s3service := services.NewS3Service(s3client, AppConfig)

	router := mux.NewRouter()
	v1Router := router.PathPrefix("/v1").Subrouter()

	v1Router.Handle("/files/upload/", handlers.NewUploadFileHandler(s3service)).Methods("POST")

	server := &http.Server{
		Addr:    ":" + AppConfig.AppPort,
		Handler: router,
	}

	log.Infof("Starting server on :%s...", AppConfig.AppPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
