package fastimagehash

/*
#cgo amd64 CFLAGS: -I${SRCDIR}/include/amd64
#cgo arm64 CFLAGS: -I${SRCDIR}/include/arm64
#cgo amd64 LDFLAGS: -L${SRCDIR}/lib/amd64
#cgo arm64 LDFLAGS: -L${SRCDIR}/lib/arm64
#cgo darwin LDFLAGS: -lscale_normal_darwin
#cgo darwin LDFLAGS: -lscale_box_darwin
#cgo darwin LDFLAGS: -lscale_linear_darwin
#cgo darwin LDFLAGS: -lscale_gauss_darwin
#cgo linux LDFLAGS: -lscale_normal_linux
#cgo linux LDFLAGS: -lscale_box_linux
#cgo linux LDFLAGS: -lscale_linear_linux
#cgo linux LDFLAGS: -lscale_gauss_linux
#cgo LDFLAGS: -ldl -lm

#include "scale_normal.h"
#include "scale_box.h"
#include "scale_linear.h"
#include "scale_gauss.h"
*/
import "C"

import (
	"fmt"
	"image"

	_ "github.com/benesch/cgosymbolizer"
)

//go:generate go run ./cmd/compile f scale_normal scale.cpp
func scaleNormal(in *image.RGBA, scaleWidth, scaleHeight int) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()

	out := image.NewRGBA(image.Rect(0, 0, scaleWidth, scaleHeight))
	outBuf, err := HalideBufferRGBA(out.Pix, scaleWidth, scaleHeight)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(outBuf)

	inBuf, err := HalideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.scale_normal(
		inBuf,
		C.int(width),
		C.int(height),
		C.int(scaleWidth),
		C.int(scaleHeight),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, fmt.Errorf("failed to scale_normal")
	}
	return out, nil
}

//go:generate go run ./cmd/compile f scale_box scale.cpp
func scaleBox(in *image.RGBA, scaleWidth, scaleHeight int) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()

	out := image.NewRGBA(image.Rect(0, 0, scaleWidth, scaleHeight))
	outBuf, err := HalideBufferRGBA(out.Pix, scaleWidth, scaleHeight)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(outBuf)

	inBuf, err := HalideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.scale_box(
		inBuf,
		C.int(width),
		C.int(height),
		C.int(scaleWidth),
		C.int(scaleHeight),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, fmt.Errorf("failed to scale_box")
	}
	return out, nil
}

//go:generate go run ./cmd/compile f scale_linear scale.cpp
func scaleLinear(in *image.RGBA, scaleWidth, scaleHeight int) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()

	out := image.NewRGBA(image.Rect(0, 0, scaleWidth, scaleHeight))
	outBuf, err := HalideBufferRGBA(out.Pix, scaleWidth, scaleHeight)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(outBuf)

	inBuf, err := HalideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.scale_linear(
		inBuf,
		C.int(width),
		C.int(height),
		C.int(scaleWidth),
		C.int(scaleHeight),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, fmt.Errorf("failed to scale_linear")
	}
	return out, nil
}

//go:generate go run ./cmd/compile f scale_gauss scale.cpp
func scaleGauss(in *image.RGBA, scaleWidth, scaleHeight int) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()

	out := image.NewRGBA(image.Rect(0, 0, scaleWidth, scaleHeight))
	outBuf, err := HalideBufferRGBA(out.Pix, scaleWidth, scaleHeight)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(outBuf)

	inBuf, err := HalideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.scale_gauss(
		inBuf,
		C.int(width),
		C.int(height),
		C.int(scaleWidth),
		C.int(scaleHeight),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, fmt.Errorf("failed to scale_gauss")
	}
	return out, nil
}
