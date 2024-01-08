package fastimagehash

import (
	"image"
	"sort"

	"github.com/pkg/errors"
)

func PerceptionHash(img *image.RGBA) (Hash64, error) {
	scaled, err := scaleNormal(img, 64, 64)
	if err != nil {
		return Hash64{}, errors.WithStack(err)
	}

	// [gray, gray, gray, alpha]
	grayRGBA, err := grayscale(scaled)
	if err != nil {
		return Hash64{}, errors.WithStack(err)
	}

	// [gray] = [gray, gray, gray, alpha]
	gray := rgbaGrayscaleToGray(grayRGBA.Pix)

	sums, err := dct2d(gray, 64, 64)
	if err != nil {
		return Hash64{}, errors.WithStack(err)
	}
	median := pickMedian(sums)
	hash := phash64x1(sums, median)
	return Hash64{hash}, nil
}

func pickMedian(src []float32) float32 {
	tmp := make([]float32, len(src))
	copy(tmp, src)

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i] < tmp[j]
	})
	half := len(src) / 2
	return tmp[half]
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
