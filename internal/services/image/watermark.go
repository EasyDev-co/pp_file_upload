package image

import (
	"github.com/h2non/bimg"
)

func (s *imageService) Watermark(original []byte) ([]byte, error) {
	options := bimg.WatermarkImage{
		Buf:     s.cfg.WatermarkBytes,
		Opacity: 0.5,
		Top:     10,
		Left:    10,
	}

	result, err := bimg.NewImage(original).WatermarkImage(options)
	if err != nil {
		return nil, err
	}

	return result, nil
}
