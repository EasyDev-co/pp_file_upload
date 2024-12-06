package api

type PhotoResponse struct {
	OriginalPhoto    string `json:"original_photo"`
	WatermarkedPhoto string `json:"watermarked_photo"`
}
