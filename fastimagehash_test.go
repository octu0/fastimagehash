package fastimagehash

import (
	"bytes"
	"image"
	"image/color"
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

	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y += 1 {
		for x := b.Min.X; x < b.Max.X; x += 1 {
			c := color.RGBAModel.Convert(img.At(x, y))
			rgba.Set(x, y, c)
		}
	}
	return rgba, nil
}

func saveImage(img *image.RGBA) (string, error) {
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
