package dto

type UploadedFilesDTO struct {
	FileURL string
}

type ProcessedFileDTO struct {
	Name               string // Имя файла
	OriginalContent    []byte
	WatermarkedContent []byte
}
