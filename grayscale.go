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

	inBuf, err := halideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(inBuf)

	out := image.NewRGBA(image.Rect(0, 0, width, height))
	outBuf, err := halideBufferRGBA(out.Pix, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(outBuf)

	ret := C.grayscale(
		inBuf,
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, errors.WithStack(ErrGrayscale)
	}
	return out, nil
}
