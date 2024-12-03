package repository

import (
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"io"
)

type S3ServiceInterface interface {
	BulkUpload(fileReaders []struct {
		Name    string
		Content io.Reader
	}) (*[]dto.UploadedFilesDTO, error)
	BulkDelete(files []string) error
}
