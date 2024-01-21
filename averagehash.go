package fastimagehash

import (
	"image"

	"github.com/pkg/errors"
)

func AverageHash(in image.Image) (Hash1024, error) {
	img, err := ConvertToRGBA(in)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}

	scaled, err := scaleNormal(img, 32, 32)
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

	avg4 := average16x16(gray)
	return avghash16x16(gray, avg4), nil
}

func average16x16(gray []byte) [4]uint32 {
	return [4]uint32{
		averageUint32(gray[0:256]),
		averageUint32(gray[256:512]),
		averageUint32(gray[512:768]),
		averageUint32(gray[768:1024]),
	}
}

func avghash16x16(gray []byte, avg4 [4]uint32) Hash1024 {
	return [16]uint64{
		avghash8x8(gray[0:64], avg4[0]),
		avghash8x8(gray[64:128], avg4[0]),
		avghash8x8(gray[128:192], avg4[0]),
		avghash8x8(gray[192:256], avg4[0]),
		avghash8x8(gray[256:320], avg4[1]),
		avghash8x8(gray[320:384], avg4[1]),
		avghash8x8(gray[384:448], avg4[1]),
		avghash8x8(gray[448:512], avg4[1]),
		avghash8x8(gray[512:576], avg4[2]),
		avghash8x8(gray[576:640], avg4[2]),
		avghash8x8(gray[640:704], avg4[2]),
		avghash8x8(gray[704:768], avg4[2]),
		avghash8x8(gray[768:832], avg4[3]),
		avghash8x8(gray[832:896], avg4[3]),
		avghash8x8(gray[896:960], avg4[3]),
		avghash8x8(gray[960:1024], avg4[3]),
	}
}

func averageUint32(values []byte) uint32 {
	sum := uint32(0)
	for _, v := range values {
		sum += uint32(v)
	}
	return sum / uint32(len(values))
}

func avghash8x8(values []byte, avg uint32) uint64 {
	index := uint64(0)
	for i := 0; i < 64; i += 1 {
		if avg < uint32(values[i]) {
			index |= 1 << (64 - i - 1)
		}
	}
	return index
}
