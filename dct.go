package fastimagehash

/*
#cgo amd64 CFLAGS: -I${SRCDIR}/include/amd64
#cgo arm64 CFLAGS: -I${SRCDIR}/include/arm64
#cgo amd64 LDFLAGS: -L${SRCDIR}/lib/amd64
#cgo arm64 LDFLAGS: -L${SRCDIR}/lib/arm64
#cgo darwin LDFLAGS: -ldct2d_darwin
#cgo linux LDFLAGS: -ldct2d_linux
#cgo LDFLAGS: -ldl -lm

#include "dct2d.h"
*/
import "C"

import (
	"github.com/pkg/errors"

	_ "github.com/benesch/cgosymbolizer"
)

var (
	ErrDCT = errors.New("failed to dct")
)

//go:generate go run ./cmd/compile f dct2d dct.cpp
func dct2d(in []byte, width, height int) ([]float32, error) {
	out := make([]float32, width*height)
	outBuf, err := halideBuffer2DFloat32(out, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(outBuf)

	inBuf, err := halideBuffer2DUint8(in, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(inBuf)

	ret := C.dct2d(
		inBuf,
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, errors.WithStack(ErrDCT)
	}
	return out, nil
}
