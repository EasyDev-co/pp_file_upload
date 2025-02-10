package image

import (
	"bytes"

	"github.com/pkg/errors"

	"image"
	"image/jpeg"
	_ "image/png"
)

const (
	jpegQuality = 20
)

// Compress сжатие изображений
func (s *imageService) Compress(input []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	if format != "jpeg" && format != "png" {
		return nil, errors.New("unsupported image format")
	}

	var buf bytes.Buffer

	opt := jpeg.Options{Quality: jpegQuality}
	err = jpeg.Encode(&buf, img, &opt)

	return buf.Bytes(), nil
}
