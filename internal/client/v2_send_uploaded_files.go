package client

import (
	"EasyDev-co/pp_file_upload/internal/consts"
	"EasyDev-co/pp_file_upload/internal/model/api"
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// V2SendUploadedFiles отправляет POST запрос для сохранения загруженных файлов
func (c *Client) V2SendUploadedFiles(ctx context.Context, sortedFiles []dto.SortedFilesDTO, kindergartenID string) error {
	endpoint := "/api/v2/photo/save_photos/"
	userId, ok := ctx.Value(consts.UserIdKey).(string)
	if !ok {
		return fmt.Errorf("value userId isn't string")
	}
	log.Infof("userId: %s", userId)

	tokenString, err := c.jwtService.GenerateJWT(userId)
	if err != nil {
		log.Errorf("Error generating token: %v", err)
		return fmt.Errorf("error generating token: %v", err)
	}

	requestBody := api.RequestBody{
		KindergartenID: kindergartenID,
		Photos:         []api.PhotoPayload{},
	}

	for _, file := range sortedFiles {
		requestBody.Photos = append(
			requestBody.Photos, api.PhotoPayload{
				OriginalPhoto:    file.OriginalContent,
				WatermarkedPhoto: file.WatermarkedContent,
				FileNumber:       file.FileNumber,
			},
		)
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	response, err := c.makeRequest(ctx, http.MethodPost, endpoint, nil, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + tokenString,
	}, bytes.NewBuffer(requestBodyJSON))

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}
