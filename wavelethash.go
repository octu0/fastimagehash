package fastimagehash

import (
	"image"

	"github.com/pkg/errors"
)

func WaveletHash(in image.Image) (Hash1024, error) {
	img, err := ConvertToRGBA(in)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}

	scaled, err := scaleNormal(img, 64, 64)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}

	// [gray, gray, gray, alpha]
	grayRGBA, err := grayscale(scaled)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}

	// [gray] = [gray, gray, gray, alpha]
	gray := rgbaGrayscaleToGray(grayRGBA.Pix)

	hi, err := haarHi(gray, 64, 64)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}
	median := pickMedian(hi)
	return wavelethash64x16(hi, median), nil
}

func wavelethash64x16(hi []float32, median float32) Hash1024 {
	return [16]uint64{
		wavelethash64x1(hi[0:64], median),
		wavelethash64x1(hi[64:128], median),
		wavelethash64x1(hi[128:192], median),
		wavelethash64x1(hi[192:256], median),
		wavelethash64x1(hi[256:320], median),
		wavelethash64x1(hi[320:384], median),
		wavelethash64x1(hi[384:448], median),
		wavelethash64x1(hi[448:512], median),
		wavelethash64x1(hi[512:576], median),
		wavelethash64x1(hi[576:640], median),
		wavelethash64x1(hi[640:704], median),
		wavelethash64x1(hi[704:768], median),
		wavelethash64x1(hi[768:832], median),
		wavelethash64x1(hi[832:896], median),
		wavelethash64x1(hi[896:960], median),
		wavelethash64x1(hi[960:1024], median),
	}
}

func wavelethash64x1(signal []float32, median float32) uint64 {
	index := uint64(0)
	for i := 0; i < 64; i += 1 {
		if median < signal[i] {
			index |= 1 << (64 - i - 1)
		}
	}
	return index
}
