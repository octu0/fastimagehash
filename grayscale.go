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
	"fmt"
	"image"

	_ "github.com/benesch/cgosymbolizer"
)

//go:generate go run ./cmd/compile f grayscale grayscale.cpp
func grayscale(in *image.RGBA) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()

	out := image.NewRGBA(image.Rect(0, 0, width, height))
	outBuf, err := HalideBufferRGBA(out.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(outBuf)

	inBuf, err := HalideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.grayscale(
		inBuf,
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, fmt.Errorf("failed to rotate90")
	}
	return out, nil
}
