package services

import "EasyDev-co/pp_file_upload/internal/model/dto"

type ImageService interface {
	Compress(input []byte) ([]byte, error)
	Watermark(original []byte, logo []byte) ([]byte, error)
	Upload(files []struct {
		Name               string
		OriginalContent    []byte
		WatermarkedContent []byte
	}) (*[]dto.UploadedFilesDTO, error)
}
