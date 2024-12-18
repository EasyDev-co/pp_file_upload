package client

import (
	"EasyDev-co/pp_file_upload/internal/model/dto"
	"context"
)

type PhotoUploadClient interface {
	V2SendUploadedFiles(ctx context.Context, sortedFiles []dto.SortedFilesDTO, kindergartenID string) error
}
