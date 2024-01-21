package fastimagehash

import (
	"image"

	"github.com/pkg/errors"
)

func DifferenceHash(in image.Image) (Hash1024, error) {
	img, err := ConvertToRGBA(in)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}

	scaled, err := scaleNormal(img, 33, 32)
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

	return differenceHash33x32(gray), nil
}

func differenceHash33x32(gray []byte) Hash1024 {
	return [16]uint64{
		diffhash33x2(gray[0:33], gray[33:66]),
		diffhash33x2(gray[66:99], gray[99:132]),
		diffhash33x2(gray[132:165], gray[165:198]),
		diffhash33x2(gray[198:231], gray[231:264]),
		diffhash33x2(gray[264:297], gray[297:330]),
		diffhash33x2(gray[330:363], gray[363:396]),
		diffhash33x2(gray[396:429], gray[429:462]),
		diffhash33x2(gray[462:495], gray[495:528]),
		diffhash33x2(gray[528:561], gray[561:594]),
		diffhash33x2(gray[594:627], gray[627:660]),
		diffhash33x2(gray[660:693], gray[693:726]),
		diffhash33x2(gray[726:759], gray[759:792]),
		diffhash33x2(gray[792:825], gray[825:858]),
		diffhash33x2(gray[858:891], gray[891:924]),
		diffhash33x2(gray[924:957], gray[957:990]),
		diffhash33x2(gray[990:1023], gray[1023:1056]),
	}
}

func diffhash33x2(a, b []byte) uint64 {
	index := uint64(0)
	pos := 0
	for i := 0; i < 32; i += 1 {
		if a[i] < a[i+1] {
			index |= 1 << (64 - pos - 1)
		}
		pos += 1
	}
	for i := 0; i < 32; i += 1 {
		if b[i] < b[i+1] {
			index |= 1 << (64 - pos - 1)
		}
		pos += 1
	}
	return index
}
