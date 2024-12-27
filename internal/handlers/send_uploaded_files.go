package handlers

import (
	"EasyDev-co/pp_file_upload/internal/client"
	"EasyDev-co/pp_file_upload/internal/config"
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"EasyDev-co/pp_file_upload/internal/response"
	"EasyDev-co/pp_file_upload/internal/services"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"sort"
)

type SendUploadedFilesHandler struct {
	imageService services.ImageService
	client       client.PhotoUploadClient
	cfg          config.Config
}

func NewSendUploadedFilesHandler(
	imageService services.ImageService,
	cfg config.Config,
	client client.PhotoUploadClient,
) *SendUploadedFilesHandler {
	return &SendUploadedFilesHandler{
		imageService: imageService,
		client:       client,
		cfg:          cfg,
	}
}

func (h *SendUploadedFilesHandler) ServeFastHTTP(ctx *fasthttp.RequestCtx) {
	kindergarten := string(ctx.FormValue("kindergarten"))
	kindergartenID := string(ctx.FormValue("kindergarten_id"))
	photoTheme := string(ctx.FormValue("photo_theme"))
	region := string(ctx.FormValue("region"))

	if kindergarten == "" || photoTheme == "" || region == "" || kindergartenID == "" {
		response.RespondError(
			ctx,
			fasthttp.StatusBadRequest,
			"Missing required query parameters: kindergarten, photo_theme, region",
		)
		return
	}

	files, err := h.getMultipartFiles(ctx)
	if err != nil {
		log.Warn("Error getting multipart files: ", err)
		response.RespondError(
			ctx,
			fasthttp.StatusBadRequest,
			err.Error(),
		)
		return
	}

	processedFiles, err := h.processFiles(files)
	if err != nil {
		log.Fatalf("Error processing files: %v", err)
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

	sort.Slice(uploadedFiles, func(i, j int) bool {
		return uploadedFiles[i].FileNumber < uploadedFiles[j].FileNumber
	})

	err = h.client.V2SendUploadedFiles(ctx, uploadedFiles, kindergartenID)
	if err != nil {
		response.RespondError(
			ctx,
			fasthttp.StatusInternalServerError,
			fmt.Sprintf("failed to send request %v", err),
		)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *SendUploadedFilesHandler) getMultipartFiles(ctx *fasthttp.RequestCtx) ([]*multipart.FileHeader, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("unable to parse form")
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files uploaded")
	}
	return files, nil
}

func (h *SendUploadedFilesHandler) processFiles(files []*multipart.FileHeader) ([]dto.ProcessedFileDTO, error) {
	results := make(chan dto.ProcessedFileDTO, len(files))
	defer close(results)

	maxGoroutines := 10
	sem := make(chan struct{}, maxGoroutines)

	fileCnt := 0
	for _, fileHeader := range files {
		fileCnt++
		sem <- struct{}{}
		go func(fileHeader *multipart.FileHeader, fileNumber int) {
			defer func() { <-sem }()

			file, err := fileHeader.Open()
			if err != nil {
				log.WithError(err).Error("Error opening fileHeader:", fileHeader.Filename)
				results <- dto.ProcessedFileDTO{}
				return
			}
			defer file.Close()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, file)
			if err != nil {
				log.WithError(err).Error("Error copy file to buffer:", fileHeader.Filename)
				results <- dto.ProcessedFileDTO{}
				return
			}

			compressedData, err := h.imageService.Compress(buf.Bytes())
			if err != nil {
				log.WithError(err).Error("Error compressing file:", fileHeader.Filename)
				results <- dto.ProcessedFileDTO{}
				return
			}

			watermarkedData, err := h.imageService.Watermark(compressedData)
			if err != nil {
				log.WithError(err).Error("Error making watermark for file:", fileHeader.Filename)
				results <- dto.ProcessedFileDTO{}
				return
			}

			log.Infof("Success: %s", fileHeader.Filename)
			results <- dto.ProcessedFileDTO{
				Name:               fileHeader.Filename,
				OriginalContent:    compressedData,
				WatermarkedContent: watermarkedData,
				FileNumber:         fileNumber,
			}
		}(fileHeader, fileCnt)
	}

	var processedFiles []dto.ProcessedFileDTO
	for i := 0; i < len(files); i++ {
		result := <-results
		processedFiles = append(processedFiles, result)
	}
	return processedFiles, nil
}
