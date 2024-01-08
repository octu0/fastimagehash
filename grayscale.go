package fastimagehash

/*
#cgo amd64 CFLAGS: -I${SRCDIR}/include/amd64
#cgo arm64 CFLAGS: -I${SRCDIR}/include/arm64
#cgo amd64 LDFLAGS: -L${SRCDIR}/lib/amd64
#cgo arm64 LDFLAGS: -L${SRCDIR}/lib/arm64
#cgo darwin LDFLAGS: -lgrayscale_darwin
#cgo linux LDFLAGS: -lgrayscale_linux
#cgo LDFLAGS: -ldl -lm

#include "grayscale.h"
*/
import "C"

import (
	"image"

	"github.com/pkg/errors"

	_ "github.com/benesch/cgosymbolizer"
)

var (
	ErrGrayscale = errors.New("failed to grayscale")
)

//go:generate go run ./cmd/compile f grayscale grayscale.cpp
func grayscale(in *image.RGBA) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()

	out := image.NewRGBA(image.Rect(0, 0, width, height))
	if err := rgbaGrayscale(out.Pix, in.Pix, width, height); err != nil {
		return nil, errors.WithStack(err)
	}
	return out, nil
}

func rgbaGrayscale(out []byte, in []byte, width, height int) error {
	inBuf, err := HalideBufferRGBA(in, width, height)
	if err != nil {
		return errors.WithStack(err)
	}
	defer HalideFreeBuffer(inBuf)

	outBuf, err := HalideBufferRGBA(out, width, height)
	if err != nil {
		return errors.WithStack(err)
	}
	defer HalideFreeBuffer(outBuf)

	ret := C.grayscale(
		inBuf,
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return errors.WithStack(ErrGrayscale)
	}
	return nil
}

func rgbaGrayscaleToGray(rgba []byte) []byte {
	out := make([]byte, 0, len(rgba)/4)
	for i := 0; i < len(rgba); i += 4 {
		out = append(out, rgba[i])
	}
	return out
}
