package fastimagehash

import (
	"image"

	"github.com/pkg/errors"
)

func PerceptionHash(in image.Image) (Hash1024, error) {
	img, err := ConvertToRGBA(in)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}

	scaled, err := scaleNormal(img, 64, 16)
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

	sums, err := dct2d(gray, 64, 16)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}
	median := pickMedian(sums)
	return phash64x16(sums, median), nil
}

func phash64x16(sums []float32, median float32) Hash1024 {
	return [16]uint64{
		phash64x1(sums[0:64], median),
		phash64x1(sums[64:128], median),
		phash64x1(sums[128:192], median),
		phash64x1(sums[192:256], median),
		phash64x1(sums[256:320], median),
		phash64x1(sums[320:384], median),
		phash64x1(sums[384:448], median),
		phash64x1(sums[448:512], median),
		phash64x1(sums[512:576], median),
		phash64x1(sums[576:640], median),
		phash64x1(sums[640:704], median),
		phash64x1(sums[704:768], median),
		phash64x1(sums[768:832], median),
		phash64x1(sums[832:896], median),
		phash64x1(sums[896:960], median),
		phash64x1(sums[960:1024], median),
	}
}

func phash64x1(sums []float32, median float32) uint64 {
	index := uint64(0)
	for i := 0; i < 64; i += 1 {
		if median < sums[i] {
			index |= 1 << (64 - i - 1)
		}
	}
	return index
}
