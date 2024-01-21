package fastimagehash

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestYUV444ToRGBA(t *testing.T) {
	src, err := png.Decode(bytes.NewReader(testPngImg))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	yuv444 := convertYCbCrModel(src, image.YCbCrSubsampleRatio444)

	outRGBA, err := ConvertToRGBA(yuv444)
	if err != nil {
		t.Errorf("%+v", err)
	}
	rgbaPath, err := saveImage(outRGBA)
	if err != nil {
		t.Errorf("%+v", err)
	}
	t.Logf("rgba(from yuv444) = %s", rgbaPath)
}

func TestYUV422ToRGBA(t *testing.T) {
	src, err := png.Decode(bytes.NewReader(testPngImg))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	yuv422 := convertYCbCrModel(src, image.YCbCrSubsampleRatio422)

	outRGBA, err := ConvertToRGBA(yuv422)
	if err != nil {
		t.Errorf("%+v", err)
	}
	rgbaPath, err := saveImage(outRGBA)
	if err != nil {
		t.Errorf("%+v", err)
	}
	t.Logf("rgba(from yuv422) = %s", rgbaPath)
}

func TestYUV420ToRGBA(t *testing.T) {
	src, err := png.Decode(bytes.NewReader(testPngImg))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	yuv420 := convertYCbCrModel(src, image.YCbCrSubsampleRatio420)

	outRGBA, err := ConvertToRGBA(yuv420)
	if err != nil {
		t.Errorf("%+v", err)
	}
	rgbaPath, err := saveImage(outRGBA)
	if err != nil {
		t.Errorf("%+v", err)
	}
	t.Logf("rgba(from yuv420) = %s", rgbaPath)
}
