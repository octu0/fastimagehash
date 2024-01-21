package fastimagehash

/*
#cgo amd64 CFLAGS: -I${SRCDIR}/include/amd64
#cgo arm64 CFLAGS: -I${SRCDIR}/include/arm64
#cgo amd64 LDFLAGS: -L${SRCDIR}/lib/amd64
#cgo arm64 LDFLAGS: -L${SRCDIR}/lib/arm64
#cgo darwin LDFLAGS: -lyuv444_to_rgba_darwin
#cgo darwin LDFLAGS: -lyuv422_to_rgba_darwin
#cgo darwin LDFLAGS: -lyuv420_to_rgba_darwin
#cgo linux LDFLAGS: -lyuv444_to_rgba_linux
#cgo linux LDFLAGS: -lyuv422_to_rgba_linux
#cgo linux LDFLAGS: -lyuv420_to_rgba_linux
#cgo LDFLAGS: -ldl -lm

#include "yuv444_to_rgba.h"
#include "yuv422_to_rgba.h"
#include "yuv420_to_rgba.h"
*/
import "C"

import (
	"image"
	"image/color"

	"github.com/pkg/errors"

	_ "github.com/benesch/cgosymbolizer"
)

var (
	ErrYUV444ToRGBA                = errors.New("failed to yuv444_to_rgba")
	ErrYUV422ToRGBA                = errors.New("failed to yuv422_to_rgba")
	ErrYUV420ToRGBA                = errors.New("failed to yuv420_to_rgba")
	ErrYUVNotSupportSubsampleRatio = errors.New("not support subsample ratio")
)

//go:generate go run ./cmd/compile f yuv444_to_rgba image.cpp
func yuv444ToRGBA(y, u, v []byte, strideY, strideU, strideV int, width, height int) (*image.RGBA, error) {
	out := make([]byte, width*height*4)

	yBuf, err := halideBufferYUV(y, strideY, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(yBuf)

	uBuf, err := halideBufferYUV(u, strideU, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(uBuf)

	vBuf, err := halideBufferYUV(v, strideV, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(vBuf)

	outBuf, err := halideBufferRGBA(out, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(outBuf)

	ret := C.yuv444_to_rgba(
		yBuf,
		uBuf,
		vBuf,
		C.int(strideY),
		C.int(strideU),
		C.int(strideV),
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, errors.WithStack(ErrYUV444ToRGBA)
	}
	return &image.RGBA{
		Pix:    out,
		Stride: width * 4,
		Rect:   image.Rect(0, 0, width, height),
	}, nil
}

//go:generate go run ./cmd/compile f yuv422_to_rgba image.cpp
func yuv422ToRGBA(y, u, v []byte, strideY, strideU, strideV int, width, height int) (*image.RGBA, error) {
	out := make([]byte, width*height*4)

	yBuf, err := halideBufferYUV(y, strideY, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(yBuf)

	uBuf, err := halideBufferYUV(u, strideU, width/2, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(uBuf)

	vBuf, err := halideBufferYUV(v, strideV, width/2, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(vBuf)

	outBuf, err := halideBufferRGBA(out, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(outBuf)

	ret := C.yuv422_to_rgba(
		yBuf,
		uBuf,
		vBuf,
		C.int(strideY),
		C.int(strideU),
		C.int(strideV),
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, errors.WithStack(ErrYUV422ToRGBA)
	}
	return &image.RGBA{
		Pix:    out,
		Stride: width * 4,
		Rect:   image.Rect(0, 0, width, height),
	}, nil
}

//go:generate go run ./cmd/compile f yuv420_to_rgba image.cpp
func yuv420ToRGBA(y, u, v []byte, strideY, strideU, strideV int, width, height int) (*image.RGBA, error) {
	out := make([]byte, width*height*4)

	yBuf, err := halideBufferYUV(y, strideY, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(yBuf)

	uBuf, err := halideBufferYUV(u, strideU, width/2, height/2)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(uBuf)

	vBuf, err := halideBufferYUV(v, strideV, width/2, height/2)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(vBuf)

	outBuf, err := halideBufferRGBA(out, width, height)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer halideBufferFree(outBuf)

	ret := C.yuv420_to_rgba(
		yBuf,
		uBuf,
		vBuf,
		C.int(strideY),
		C.int(strideU),
		C.int(strideV),
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, errors.WithStack(ErrYUV420ToRGBA)
	}
	return &image.RGBA{
		Pix:    out,
		Stride: width * 4,
		Rect:   image.Rect(0, 0, width, height),
	}, nil
}

func YCbCrToRGBA(ycbcr *image.YCbCr) (*image.RGBA, error) {
	width, height := ycbcr.Rect.Dx(), ycbcr.Rect.Dy()

	switch ycbcr.SubsampleRatio {
	case image.YCbCrSubsampleRatio444:
		return yuv444ToRGBA(
			ycbcr.Y,
			ycbcr.Cb,
			ycbcr.Cr,
			ycbcr.YStride,
			ycbcr.CStride,
			ycbcr.CStride,
			width,
			height,
		)
	case image.YCbCrSubsampleRatio422:
		return yuv422ToRGBA(
			ycbcr.Y,
			ycbcr.Cb,
			ycbcr.Cr,
			ycbcr.YStride,
			ycbcr.CStride,
			ycbcr.CStride,
			width,
			height,
		)
	case image.YCbCrSubsampleRatio420:
		return yuv420ToRGBA(
			ycbcr.Y,
			ycbcr.Cb,
			ycbcr.Cr,
			ycbcr.YStride,
			ycbcr.CStride,
			ycbcr.CStride,
			width,
			height,
		)
	default:
		return nil, errors.Wrapf(ErrYUVNotSupportSubsampleRatio, "ratio=%d", ycbcr.SubsampleRatio)
	}
}

func ConvertToRGBA(img image.Image) (*image.RGBA, error) {
	switch img.(type) {
	case *image.YCbCr:
		return YCbCrToRGBA(img.(*image.YCbCr))
	case *image.RGBA:
		return img.(*image.RGBA), nil
	case *image.NRGBA:
		// todo: calc non-alpha-premultiplied
		v := img.(*image.NRGBA)
		return &image.RGBA{
			Pix:    v.Pix,
			Stride: v.Stride,
			Rect:   v.Rect,
		}, nil
	default:
		return convertRGBAModel(img), nil
	}
}

func convertYCbCrModel(img image.Image, subsample image.YCbCrSubsampleRatio) *image.YCbCr {
	b := img.Bounds()
	yuv := image.NewYCbCr(img.Bounds(), subsample)
	for y := 0; y < b.Dy(); y += 1 {
		for x := 0; x < b.Dx(); x += 1 {
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			yy, uu, vv := color.RGBToYCbCr(rgba.R, rgba.G, rgba.B)

			cy := yuv.YOffset(x, y)
			ci := yuv.COffset(x, y)
			yuv.Y[cy] = yy
			yuv.Cb[ci] = uu
			yuv.Cr[ci] = vv
		}
	}
	return yuv
}

func convertRGBAModel(img image.Image) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y += 1 {
		for x := b.Min.X; x < b.Max.X; x += 1 {
			c := color.RGBAModel.Convert(img.At(x, y))
			rgba.Set(x, y, c)
		}
	}
	return rgba
}
