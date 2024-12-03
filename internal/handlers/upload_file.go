package handlers

import (
	"EasyDev-co/pp_file_upload/internal/config"
	"EasyDev-co/pp_file_upload/internal/services"
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"os"
)

type UploadFileHandler struct {
	imageService services.ImageService
	cfg          config.Config
}

func NewUploadFileHandler(imageService services.ImageService, cfg config.Config) *UploadFileHandler {
	return &UploadFileHandler{
		imageService: imageService,
		cfg:          cfg,
	}
}

func (h *UploadFileHandler) ServeFastHTTP(ctx *fasthttp.RequestCtx) {
	logoBytes, err := os.ReadFile(h.cfg.WatermarkPath)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Error reading watermark file")
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Unable to parse form")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("No files uploaded")
		return
	}

	type ProcessedFile struct {
		Name               string
		OriginalContent    []byte
		WatermarkedContent []byte
		Error              error
	}

	results := make(chan ProcessedFile, len(files))
	defer close(results)

	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)

	for _, fileHeader := range files {
		sem <- struct{}{}
		go func(fileHeader *multipart.FileHeader) {
			defer func() { <-sem }()

			file, err := fileHeader.Open()
			if err != nil {
				results <- ProcessedFile{Error: fmt.Errorf("failed to open file %s: %v", fileHeader.Filename, err)}
				return
			}
			defer file.Close()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, file)
			if err != nil {
				results <- ProcessedFile{Error: fmt.Errorf("failed to read file %s: %v", fileHeader.Filename, err)}
				return
			}

			compressedData, err := h.imageService.Compress(buf.Bytes())
			if err != nil {
				results <- ProcessedFile{Error: fmt.Errorf("failed to compress file %s: %v", fileHeader.Filename, err)}
				return
			}

			watermarkedData, err := h.imageService.Watermark(compressedData, logoBytes)
			if err != nil {
				results <- ProcessedFile{Error: fmt.Errorf("failed to add watermark to file %s: %v", fileHeader.Filename, err)}
				return
			}

			results <- ProcessedFile{
				Name:               fileHeader.Filename,
				OriginalContent:    compressedData,
				WatermarkedContent: watermarkedData,
			}
		}(fileHeader)
	}

	var processedFiles []ProcessedFile

	// Сбор результатов
	for i := 0; i < len(files); i++ {
		result := <-results
		if result.Error != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBodyString(fmt.Sprintf("Error processing file: %v", result.Error))
			return
		}
		processedFiles = append(processedFiles, result)
	}

	var uploadFiles []struct {
		Name               string
		OriginalContent    []byte
		WatermarkedContent []byte
	}
	for _, file := range processedFiles {
		uploadFiles = append(uploadFiles, struct {
			Name               string
			OriginalContent    []byte
			WatermarkedContent []byte
		}{
			Name:               file.Name,
			OriginalContent:    file.OriginalContent,
			WatermarkedContent: file.WatermarkedContent,
		})
	}

	uploadedFiles, err := h.imageService.Upload(uploadFiles)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(fmt.Sprintf("Failed to upload files to S3: %v", err))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(fmt.Sprintf("Successfully uploaded %d files", len(*uploadedFiles)))
}
