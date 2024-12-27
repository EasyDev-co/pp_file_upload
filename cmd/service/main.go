package main

import (
	"EasyDev-co/pp_file_upload/internal/services/jwt"
	"crypto/tls"
	"fmt"

	"EasyDev-co/pp_file_upload/internal/client"
	"EasyDev-co/pp_file_upload/internal/config"
	"EasyDev-co/pp_file_upload/internal/db"
	"EasyDev-co/pp_file_upload/internal/handlers"
	"EasyDev-co/pp_file_upload/internal/middleware"
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

	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}
	s3client, err := db.NewS3Client(appConfig)
	if err != nil {
		log.Fatalf("Error creating S3 client: %v", err)
		return
	}

	s3repository := s3.NewS3Repository(s3client, appConfig)
	imageService := image.NewImageService(s3repository, appConfig)
	jwtService := jwt.NewJWTService(appConfig.SigningKey)
	mdlware := middleware.NewMiddleware(jwtService, appConfig.AllowedOrigins)

	apiClient := client.NewClient(appConfig.BaseURL, appConfig.RequestTimeout, jwtService)

	r := router.New()

	r.POST(
		"/v1/files/upload/",
		handlers.NewUploadFileHandler(imageService, appConfig).ServeFastHTTP,
	)
	r.POST("/v1/files/send_uploaded/",
		handlers.NewSendUploadedFilesHandler(imageService, appConfig, apiClient).ServeFastHTTP,
	)

	handler := mdlware.Timer(mdlware.CORS(mdlware.JWT(r.Handler)))

	switch appConfig.AppEnv {
	case config.Dev:
		server := &fasthttp.Server{
			Handler:            handler,
			MaxRequestBodySize: int(appConfig.MaxUploadSize),
		}

		log.Infof("Starting HTTP server on :%s...", appConfig.AppPort)
		if err := server.ListenAndServe(fmt.Sprintf(":%s", appConfig.AppPort)); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	case config.Prod:
		cert, err := tls.LoadX509KeyPair(appConfig.CertFile, appConfig.KeyFile)
		if err != nil {
			log.Fatalf("Error loading SSL certificate and key: %v", err)
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		}

		server := &fasthttp.Server{
			Handler:            handler,
			MaxRequestBodySize: int(appConfig.MaxUploadSize),
			TLSConfig:          tlsConfig,
		}

		log.Infof("Starting HTTPS server on :%s...", appConfig.AppPort)
		if err := server.ListenAndServeTLS(
			fmt.Sprintf(":%s", appConfig.AppPort),
			appConfig.CertFile,
			appConfig.KeyFile,
		); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	default:
		log.Fatalf("Unsupported environment: %s", appConfig.AppEnv)
	}
}
