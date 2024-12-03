package repository

import (
	"EasyDev-co/pp_file_upload/internal/dto"
	"io"
)

type S3ServiceInterface interface {
	BulkUpload(fileReaders []struct {
		Name    string
		Content io.Reader
	}) (*dto.UploadResponse, error)
	BulkDelete(files []string) error
}
