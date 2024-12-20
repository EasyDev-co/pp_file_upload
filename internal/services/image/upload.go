package image

import (
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

func (s *imageService) Upload(
	files []dto.ProcessedFileDTO,
	kindergarten, photoTheme,
	region string,
) (*[]dto.SortedFilesDTO, error) {
	var sortedFiles []dto.SortedFilesDTO
	results := make(chan dto.SortedFilesDTO, len(files))
	errors := make(chan error, len(files))
	var wg sync.WaitGroup

	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)

	filesPtrs := make([]*dto.ProcessedFileDTO, 0, len(files))
	for i, _ := range files {
		filesPtrs = append(filesPtrs, &files[i])
	}

	for _, file := range filesPtrs {
		wg.Add(1)
		sem <- struct{}{}

		go func(*dto.ProcessedFileDTO) {
			defer wg.Done()
			defer func() { <-sem }()

			name := fmt.Sprintf("%s/%s/%s/%s", region, kindergarten, photoTheme, file.Name)

			log.WithFields(log.Fields{
				"filename":  file.Name,
				"full_name": name,
			}).Info("file-to-upload info")

			originalURL, err := s.s3service.UploadFile(
				name,
				bytes.NewReader(file.OriginalContent),
			)
			if err != nil {
				errors <- fmt.Errorf("failed to upload original file %s: %w", file.Name, err)
				return
			}

			watermarkedURL, err := s.s3service.UploadFile(
				fmt.Sprintf("%s/%s/%s/watermarked/%s", region, kindergarten, photoTheme, file.Name),
				bytes.NewReader(file.WatermarkedContent),
			)
			if err != nil {
				errors <- fmt.Errorf("failed to upload watermarked file %s: %w", file.Name, err)
				return
			}

			results <- dto.SortedFilesDTO{
				OriginalContent:    originalURL,
				WatermarkedContent: watermarkedURL,
			}
		}(file)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	var uploadErr error
	go func() {
		for err := range errors {
			uploadErr = err
		}
	}()

	for res := range results {
		sortedFiles = append(sortedFiles, res)
	}

	if uploadErr != nil {
		return nil, uploadErr
	}

	return &sortedFiles, nil
}
