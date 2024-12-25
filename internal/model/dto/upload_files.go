package dto

type UploadedFilesDTO struct {
	FileURL string
}

type SortedFilesDTO struct {
	OriginalContent    string
	WatermarkedContent string
	FileNumber         int
}

type ProcessedFileDTO struct {
	Name               string // Имя файла
	OriginalContent    []byte
	WatermarkedContent []byte
	FileNumber         int
}
