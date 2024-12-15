package image

import (
	"fmt"
	"github.com/h2non/bimg"
)

const (
	opacity      = 1.5
	shiftPixels  = 10
	logoSizeCoef = 0.75
)

func (s *imageService) Watermark(original []byte) ([]byte, error) {
	mainSize, err := bimg.NewImage(original).Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get main image size: %w", err)
	}

	logoHeight := int(float64(mainSize.Height) * logoSizeCoef)

	resizedLogo, err := bimg.NewImage(s.cfg.WatermarkBytes).ForceResize(0, logoHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to resize logo: %w", err)
	}

	logoSize, err := bimg.NewImage(resizedLogo).Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get logo size: %w", err)
	}

	centerTop := (mainSize.Height - logoSize.Height) / 2

	leftWatermark := bimg.WatermarkImage{
		Buf:     resizedLogo,
		Opacity: opacity,
		Top:     centerTop,
		Left:    shiftPixels,
	}

	leftImage, err := bimg.NewImage(original).WatermarkImage(leftWatermark)
	if err != nil {
		return nil, fmt.Errorf("failed to apply left watermark: %w", err)
	}

	rightWatermark := bimg.WatermarkImage{
		Buf:     resizedLogo,
		Opacity: opacity,
		Top:     centerTop,
		Left:    mainSize.Width - logoSize.Width - shiftPixels,
	}

	resultImage, err := bimg.NewImage(leftImage).WatermarkImage(rightWatermark)
	if err != nil {
		return nil, fmt.Errorf("failed to apply right watermark: %w", err)
	}

	return resultImage, nil
}
