package fastimagehash

import (
	"bytes"
	"image"
	"image/png"
	"os"

	_ "embed"
)

var (
	//go:embed testdata/src.png
	testPngImg []byte
	//go:embed testdata/blend.png
	testPngBlendImg []byte
	//go:embed testdata/blur.png
	testPngBlurImg []byte
	//go:embed testdata/catA_1.png
	testPngCatA1Img []byte
	//go:embed testdata/catA_2.png
	testPngCatA2Img []byte
)

func pngToRGBA(data []byte) (*image.RGBA, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if i, ok := img.(*image.RGBA); ok {
		return i, nil
	}

	return convertRGBAModel(img), nil
}

func saveImage(img image.Image) (string, error) {
	out, err := os.CreateTemp("/tmp", "out*.png")
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return "", err
	}
	return out.Name(), nil
}
