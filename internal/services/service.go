package services

import "EasyDev-co/pp_file_upload/internal/model/dto"

type ImageService interface {
	Compress(input []byte) ([]byte, error)
	Watermark(original []byte) ([]byte, error)
	Upload(files []dto.ProcessedFileDTO) (*[]dto.UploadedFilesDTO, error)
}
