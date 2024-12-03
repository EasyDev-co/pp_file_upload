package image

import (
	"EasyDev-co/pp_file_upload/internal/config"
	"EasyDev-co/pp_file_upload/internal/repository"
	"EasyDev-co/pp_file_upload/internal/services"
)

type imageService struct {
	s3service repository.S3ServiceInterface
	cfg       config.Config
}

func NewImageService(s3service repository.S3ServiceInterface, cfg config.Config) services.ImageService {
	return &imageService{s3service: s3service, cfg: cfg}
}
