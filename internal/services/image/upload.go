package image

import (
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"bytes"
	"fmt"
	"io"
)

func (s *imageService) Upload(files []dto.ProcessedFileDTO) (*[]dto.UploadedFilesDTO, error) {
	var fileReaders []struct {
		Name    string
		Content io.Reader
	}

	for _, file := range files {
		fileReaders = append(fileReaders, struct {
			Name    string
			Content io.Reader
		}{
			Name:    fmt.Sprintf("original/%s", file.Name),
			Content: bytes.NewReader(file.OriginalContent),
		})

		fileReaders = append(fileReaders, struct {
			Name    string
			Content io.Reader
		}{
			Name:    fmt.Sprintf("watermarked/%s", file.Name),
			Content: bytes.NewReader(file.WatermarkedContent),
		})
	}

	uploadedFiles, err := s.s3service.BulkUpload(fileReaders)
	if err != nil {
		return nil, fmt.Errorf("failed to upload files to S3: %w", err)
	}

	return uploadedFiles, nil
}
