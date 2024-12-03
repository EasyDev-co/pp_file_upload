package handlers

import (
	"EasyDev-co/pp_file_upload/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type UploadFileHandler struct {
	s3service repository.S3ServiceInterface
}

func NewUploadFileHandler(s3service repository.S3ServiceInterface) *UploadFileHandler {
	return &UploadFileHandler{
		s3service: s3service,
	}
}

func (uh *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("Received request to upload_file")

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Error("Failed to parse form:", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		log.Error("No files uploaded")
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	log.Info("Decoding request body")

	var fileReaders []struct {
		Name    string
		Content io.Reader
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Error("Failed to open file %s: %v", fileHeader.Filename, err)
			http.Error(w, fmt.Sprintf("Error opening file %s: %v", fileHeader.Filename, err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileReaders = append(fileReaders, struct {
			Name    string
			Content io.Reader
		}{
			Name:    fileHeader.Filename,
			Content: file,
		})
	}

	response, err := uh.s3service.BulkUpload(fileReaders)
	if err != nil {
		log.Printf("Error during file upload: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Println("File upload completed successfully")
}
