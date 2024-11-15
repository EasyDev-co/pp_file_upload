package dto

type PhotoDTO struct {
	Files []string `json:"files"`
}

type UploadResponse struct {
	FileURLs []string `json:"file_urls"`
}
