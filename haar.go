package fastimagehash

/*
#cgo amd64 CFLAGS: -I${SRCDIR}/include/amd64
#cgo arm64 CFLAGS: -I${SRCDIR}/include/arm64
#cgo amd64 LDFLAGS: -L${SRCDIR}/lib/amd64
#cgo arm64 LDFLAGS: -L${SRCDIR}/lib/arm64
#cgo darwin LDFLAGS: -lhaar_x_darwin
#cgo darwin LDFLAGS: -lhaar_y_darwin
#cgo darwin LDFLAGS: -lhaar_darwin
#cgo darwin LDFLAGS: -lhaar_hi_darwin
#cgo linux LDFLAGS: -lhaar_x_linux
#cgo linux LDFLAGS: -lhaar_y_linux
#cgo linux LDFLAGS: -lhaar_linux
#cgo linux LDFLAGS: -lhaar_hi_linux
#cgo LDFLAGS: -ldl -lm

#include "haar_x.h"
#include "haar_y.h"
#include "haar.h"
#include "haar_hi.h"
*/
import "C"

import (
	"math"

	"github.com/pkg/errors"

	_ "github.com/benesch/cgosymbolizer"
)

var (
	ErrHaarX    = errors.New("failed to haar_x")
	ErrHaarY    = errors.New("failed to haar_y")
	ErrHaarXY   = errors.New("failed to haar_xy")
	ErrHaarXYHi = errors.New("failed to haar_xy_hi")
)

//go:generate go run ./cmd/compile p haar_x haar.cpp
func haarX(in []byte, width, height int) ([]float32, []float32, error) {
	lo := make([]float32, (width/2)*height)
	hi := make([]float32, (width/2)*height)

	loBuf, err := HalideBuffer2DFloat32(lo, width/2, height)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(loBuf)

	hiBuf, err := HalideBuffer2DFloat32(hi, width/2, height)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(hiBuf)

	inBuf, err := HalideBuffer2DUint8(in, width, height)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.haar_x(
		inBuf,
		C.int(width),
		C.int(height),
		loBuf,
		hiBuf,
	)
	if ret != C.int(0) {
		return nil, nil, errors.WithStack(ErrHaarX)
	}
	return lo, hi, nil
}

//go:generate go run ./cmd/compile p haar_y haar.cpp
func haarY(in []byte, width, height int) ([]float32, []float32, error) {
	lo := make([]float32, width*(height/2))
	hi := make([]float32, width*(height/2))

	loBuf, err := HalideBuffer2DFloat32(lo, width, height/2)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(loBuf)

	hiBuf, err := HalideBuffer2DFloat32(hi, width, height/2)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(hiBuf)

	inBuf, err := HalideBuffer2DUint8(in, width, height)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.haar_y(
		inBuf,
		C.int(width),
		C.int(height),
		loBuf,
		hiBuf,
	)
	if ret != C.int(0) {
		return nil, nil, errors.WithStack(ErrHaarY)
	}
	return lo, hi, nil
}

//go:generate go run ./cmd/compile p haar haar.cpp
func haar(in []byte, width, height int) ([]float32, []float32, error) {
	lo := make([]float32, (width*height)/2)
	hi := make([]float32, (width*height)/2)

	loBuf, err := HalideBuffer2DFloat32(lo, width/2, height/2)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(loBuf)

	hiBuf, err := HalideBuffer2DFloat32(hi, width/2, height/2)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(hiBuf)

	inBuf, err := HalideBuffer2DUint8(in, width, height)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.haar(
		inBuf,
		C.int(width),
		C.int(height),
		loBuf,
		hiBuf,
	)
	if ret != C.int(0) {
		return nil, nil, errors.WithStack(ErrHaarXY)
	}
	return lo, hi, nil
}

//go:generate go run ./cmd/compile f haar_hi haar.cpp
func haarHi(in []byte, width, height int) ([]float32, error) {
	hi := make([]float32, (width*height)/2)

	hiBuf, err := HalideBuffer2DFloat32(hi, width/2, height/2)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(hiBuf)

	inBuf, err := HalideBuffer2DUint8(in, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.haar_hi(
		inBuf,
		C.int(width),
		C.int(height),
		hiBuf,
	)
	if ret != C.int(0) {
		return nil, errors.WithStack(ErrHaarXYHi)
	}
	return hi, nil
}

var (
	sqrt2 = float32(math.Sqrt2)
)

func nativeHaar1D(signal []float32) (low, high []float32) {
	N := len(signal)
	lo, hi := make([]float32, N/2), make([]float32, N/2)

	for i := 0; i < N/2; i += 1 {
		lo[i] = (signal[2*i] + signal[2*i+1]) / sqrt2
		hi[i] = (signal[2*i] - signal[2*i+1]) / sqrt2
	}
	return lo, hi
}

func nativeHaar2D(signal [][]float32) (low, high [][]float32) {
	W := len(signal[0])
	H := len(signal)

	loX := make([][]float32, H)
	hiX := make([][]float32, H)
	for h := 0; h < H; h += 1 {
		loX[h] = make([]float32, W/2)
		hiX[h] = make([]float32, W/2)
		for w := 0; w < W/2; w += 1 {
			loX[h][w] = (signal[h][2*w] + signal[h][2*w+1]) / sqrt2
			hiX[h][w] = (signal[h][2*w] - signal[h][2*w+1]) / sqrt2
		}
	}

	loXY := make([][]float32, H/2)
	hiXY := make([][]float32, H/2)
	for h := 0; h < H/2; h += 1 {
		loXY[h] = make([]float32, W/2)
		hiXY[h] = make([]float32, W/2)
		for w := 0; w < W/2; w += 1 {
			loXY[h][w] = (loX[2*h][w] + loX[2*h+1][w]) / sqrt2
			hiXY[h][w] = (hiX[2*h][w] - hiX[2*h+1][w]) / sqrt2
		}
	}
	return loXY, hiXY
}
