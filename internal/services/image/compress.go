package image

import (
	"bytes"

	"github.com/pkg/errors"

	"image"
	"image/jpeg"
	"image/png"
)

// Compress сжатие изображений
func (s *imageService) Compress(input []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	switch format {
	case "jpeg":
		opt := jpeg.Options{Quality: 75}
		err = jpeg.Encode(&buf, img, &opt)
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&buf, img)
	default:
		return nil, errors.Wrap(err, "unsupported image format")
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
