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
	return wavelethash32x32(hi, median), nil
}

func wavelethash32x32(hi []float32, median float32) Hash1024 {
	return [16]uint64{
		wavelethash32x2(hi[0:32], hi[32:64], median),
		wavelethash32x2(hi[64:96], hi[96:128], median),
		wavelethash32x2(hi[128:160], hi[160:192], median),
		wavelethash32x2(hi[192:224], hi[224:256], median),
		wavelethash32x2(hi[256:288], hi[288:320], median),
		wavelethash32x2(hi[320:352], hi[352:384], median),
		wavelethash32x2(hi[384:416], hi[416:448], median),
		wavelethash32x2(hi[448:480], hi[480:512], median),
		wavelethash32x2(hi[512:544], hi[544:576], median),
		wavelethash32x2(hi[576:608], hi[608:640], median),
		wavelethash32x2(hi[640:672], hi[672:704], median),
		wavelethash32x2(hi[704:736], hi[736:768], median),
		wavelethash32x2(hi[768:800], hi[800:832], median),
		wavelethash32x2(hi[832:864], hi[864:896], median),
		wavelethash32x2(hi[896:928], hi[928:960], median),
		wavelethash32x2(hi[960:992], hi[992:1024], median),
	}
}

func wavelethash32x2(a, b []float32, median float32) uint64 {
	index := uint64(0)
	pos := 0
	for i := 0; i < 32; i += 1 {
		if median < a[i] {
			index |= 1 << (64 - pos - 1)
		}
		pos += 1
	}
	for i := 0; i < 32; i += 1 {
		if median < b[i] {
			index |= 1 << (64 - pos - 1)
		}
		pos += 1
	}
	return index
}
