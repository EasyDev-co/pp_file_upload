package dto

type UploadedFilesDTO struct {
	FileURL string
}

type SortedFilesDTO struct {
	OriginalContent    string
	WatermarkedContent string
}

type ProcessedFileDTO struct {
	Name               string // Имя файла
	OriginalContent    *[]byte
	WatermarkedContent *[]byte
}
