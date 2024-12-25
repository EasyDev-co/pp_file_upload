package api

type PhotoResponse struct {
	OriginalPhoto    string `json:"original_photo"`
	WatermarkedPhoto string `json:"watermarked_photo"`
}

type RequestBody struct {
	KindergartenID string         `json:"kindergarten_id"`
	Photos         []PhotoPayload `json:"photos"`
}

// PhotoPayload структура для передачи данных о фотографиях
type PhotoPayload struct {
	OriginalPhoto    string `json:"original_photo"`
	WatermarkedPhoto string `json:"watermarked_photo"`
	FileNumber       int    `json:"file_number"`
}
