package services

import (
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"github.com/golang-jwt/jwt/v5"
)

type ImageService interface {
	Compress(input []byte) ([]byte, error)
	Watermark(original []byte) ([]byte, error)
	Upload(files []dto.ProcessedFileDTO, kindergarten, photoTheme, region string) ([]dto.SortedFilesDTO, error)
}

type JWTService interface {
	GenerateJWT(userId string) (string, error)
	ParseJWT(tokenString string) (*jwt.Token, error)
}
