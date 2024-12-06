package handlers

import (
	"EasyDev-co/pp_file_upload/internal/model/api"
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"EasyDev-co/pp_file_upload/internal/response"
	"EasyDev-co/pp_file_upload/internal/services"
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
)

type UploadFileHandler struct {
	imageService services.ImageService
}

func NewUploadFileHandler(imageService services.ImageService) *UploadFileHandler {
	return &UploadFileHandler{
		imageService: imageService,
	}
}

func (h *UploadFileHandler) ServeFastHTTP(ctx *fasthttp.RequestCtx) {
	kindergarten := string(ctx.QueryArgs().Peek("kindergarten"))
	photoTheme := string(ctx.QueryArgs().Peek("photo_theme"))
	region := string(ctx.QueryArgs().Peek("region"))

	if kindergarten == "" || photoTheme == "" || region == "" {
		response.RespondError(
			ctx,
			fasthttp.StatusBadRequest,
			"Missing required query parameters: kindergarten, photo_theme, region",
		)
		return
	}

	files, err := h.getMultipartFiles(ctx)
	if err != nil {
		response.RespondError(
			ctx,
			fasthttp.StatusBadRequest,
			err.Error(),
		)
		return
	}

	processedFiles, err := h.processFiles(files)
	if err != nil {
		response.RespondError(
			ctx,
			fasthttp.StatusInternalServerError,
			fmt.Sprintf("Error processing file: %v", err),
		)
		return
	}

	uploadedFiles, err := h.imageService.Upload(processedFiles, kindergarten, photoTheme, region)
	if err != nil {
		response.RespondError(
			ctx,
			fasthttp.StatusInternalServerError,
			fmt.Sprintf("Failed to upload files to S3: %v", err),
		)
		return
	}

	var photoResponses []api.PhotoResponse
	for _, file := range *uploadedFiles {
		photoResponses = append(photoResponses, api.PhotoResponse{
			OriginalPhoto:    file.OriginalContent,
			WatermarkedPhoto: file.WatermarkedContent,
		})
	}

	response.RespondSuccessJSON(ctx, photoResponses)
}

func (h *UploadFileHandler) getMultipartFiles(ctx *fasthttp.RequestCtx) ([]*multipart.FileHeader, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("unable to parse form")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files uploaded")
	}
	return files, nil
}

func (h *UploadFileHandler) processFiles(files []*multipart.FileHeader) ([]dto.ProcessedFileDTO, error) {
	results := make(chan dto.ProcessedFileDTO, len(files))
	defer close(results)

	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)

	for _, fileHeader := range files {
		sem <- struct{}{}
		go func(fileHeader *multipart.FileHeader) {
			defer func() { <-sem }()

			file, err := fileHeader.Open()
			if err != nil {
				results <- dto.ProcessedFileDTO{}
				return
			}
			defer file.Close()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, file)
			if err != nil {
				results <- dto.ProcessedFileDTO{}
				return
			}

			compressedData, err := h.imageService.Compress(buf.Bytes())
			if err != nil {
				results <- dto.ProcessedFileDTO{}
				return
			}

			watermarkedData, err := h.imageService.Watermark(compressedData)
			if err != nil {
				results <- dto.ProcessedFileDTO{}
				return
			}

			results <- dto.ProcessedFileDTO{
				Name:               fileHeader.Filename,
				OriginalContent:    compressedData,
				WatermarkedContent: watermarkedData,
			}
		}(fileHeader)
	}

	var processedFiles []dto.ProcessedFileDTO
	for i := 0; i < len(files); i++ {
		result := <-results
		processedFiles = append(processedFiles, result)
	}
	return processedFiles, nil
}
