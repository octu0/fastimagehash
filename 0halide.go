package fastimagehash

//go:generate go run ./cmd/download

/*
#cgo amd64 CFLAGS: -I${SRCDIR}/include/amd64
#cgo arm64 CFLAGS: -I${SRCDIR}/include/arm64
#cgo amd64 LDFLAGS: -L${SRCDIR}/lib/amd64
#cgo arm64 LDFLAGS: -L${SRCDIR}/lib/arm64
#cgo darwin LDFLAGS: -lruntime_darwin
#cgo linux LDFLAGS: -lruntime_linux
#cgo LDFLAGS: -ldl -lm

#include <stdint.h>
#include <stdlib.h>
#include "runtime.h"

const struct halide_type_t halide_int16_t = { halide_type_int,  16, 1 };
const struct halide_type_t halide_int32_t = { halide_type_int, 32, 1 };
const struct halide_type_t halide_uint8_t = { halide_type_uint,  8, 1 };
const struct halide_type_t halide_uint16_t = { halide_type_uint, 16, 1 };
const struct halide_type_t halide_float_t = { halide_type_float, 32, 1 };
const struct halide_type_t halide_double_t = { halide_type_float, 64, 1 };

static void free_halide_buffer(halide_buffer_t *buf) {
  if (NULL != buf) {
    free(buf->dim);
  }
  free(buf);
}

static void init_rgba_dim(halide_dimension_t *dim, int32_t width, int32_t height) {
  // width
  dim[0].min = 0;
  dim[0].extent = width;
  dim[0].stride = 4;
  dim[0].flags = 0;

  // height
  dim[1].min = 0;
  dim[1].extent = height;
  dim[1].stride = width * 4;
  dim[1].flags = 0;

  // channel
  dim[2].min = 0;
  dim[2].extent = 4;
  dim[2].stride = 1;
  dim[2].flags = 0;
}

static void init_2d_dim(halide_dimension_t *dim, int32_t width, int32_t height) {
  // width
  dim[0].min = 0;
  dim[0].extent = width;
  dim[0].stride = 1;
  dim[0].flags = 0;

  // height
  dim[1].min = 0;
  dim[1].extent = height;
  dim[1].stride = width;
  dim[1].flags = 0;
}

static halide_buffer_t *create_buffer(void *data, halide_dimension_t *dim, int dimensions, struct halide_type_t halide_type) {
  halide_buffer_t *buffer = (halide_buffer_t *) malloc(sizeof(halide_buffer_t));
  if(buffer == NULL) {
    return NULL;
  }
  memset(buffer, 0, sizeof(halide_buffer_t));

  buffer->dimensions = dimensions;
  buffer->dim = dim;
  buffer->device = 0;
  buffer->device_interface = NULL;
  buffer->host = data;
  buffer->flags = halide_buffer_flag_host_dirty;
  buffer->type = halide_type;
  return buffer;
}

static halide_buffer_t *create_halide_buffer_rgba(uint8_t *data, int width, int height) {
  int dimensions = 3;
  halide_dimension_t *dim = (halide_dimension_t *) malloc(dimensions * sizeof(halide_dimension_t));
  if(NULL == dim) {
    return NULL;
  }
  memset(dim, 0, dimensions * sizeof(halide_dimension_t));
  init_rgba_dim(dim, width, height);

  halide_buffer_t *buf = create_buffer(data, dim, dimensions, halide_uint8_t);
  if(NULL == buf) {
    free(dim);
    return NULL;
  }
  return buf;
}

static halide_buffer_t *create_halide_buffer_2d_uint8(uint8_t *data, int width, int height) {
  int dimensions = 2;
  halide_dimension_t *dim = (halide_dimension_t *) malloc(dimensions * sizeof(halide_dimension_t));
  if(NULL == dim) {
    return NULL;
  }
  memset(dim, 0, dimensions * sizeof(halide_dimension_t));
  init_2d_dim(dim, width, height);

  halide_buffer_t *buf = create_buffer(data, dim, dimensions, halide_uint8_t);
  if(NULL == buf) {
    free(dim);
    return NULL;
  }
  return buf;
}

static halide_buffer_t *create_halide_buffer_2d_float(float *data, int width, int height) {
  int dimensions = 2;
  halide_dimension_t *dim = (halide_dimension_t *) malloc(dimensions * sizeof(halide_dimension_t));
  if(NULL == dim) {
    return NULL;
  }
  memset(dim, 0, dimensions * sizeof(halide_dimension_t));
  init_2d_dim(dim, width, height);

  halide_buffer_t *buf = create_buffer(data, dim, dimensions, halide_float_t);
  if(NULL == buf) {
    free(dim);
    return NULL;
  }
  return buf;
}
*/
import "C"

import (
	"unsafe"

	"github.com/pkg/errors"

	_ "github.com/benesch/cgosymbolizer"
)

var (
	ErrHalideBufferRGBA      = errors.New("failed to create_halide_buffer_rgba")
	ErrHalideBuffer2DUint8   = errors.New("failed to create_halide_buffer_2d_uint8")
	ErrHalideBuffer2DFloat32 = errors.New("failed to create_halide_buffer_2d_float")
)

func HalideBufferRGBA(data []byte, width, height int) (*C.halide_buffer_t, error) {
	buf := unsafe.Pointer(C.create_halide_buffer_rgba(
		(*C.uint8_t)(unsafe.Pointer(&data[0])),
		C.int(width),
		C.int(height),
	))
	if buf == nil {
		return nil, errors.WithStack(ErrHalideBufferRGBA)
	}
	return (*C.halide_buffer_t)(buf), nil
}

func HalideBuffer2DUint8(data []uint8, width, height int) (*C.halide_buffer_t, error) {
	buf := unsafe.Pointer(C.create_halide_buffer_2d_uint8(
		(*C.uint8_t)(unsafe.Pointer(&data[0])),
		C.int(width),
		C.int(height),
	))
	if buf == nil {
		return nil, errors.WithStack(ErrHalideBuffer2DUint8)
	}
	return (*C.halide_buffer_t)(buf), nil
}

func HalideBuffer2DFloat32(data []float32, width, height int) (*C.halide_buffer_t, error) {
	buf := unsafe.Pointer(C.create_halide_buffer_2d_float(
		(*C.float)(unsafe.Pointer(&data[0])),
		C.int(width),
		C.int(height),
	))
	if buf == nil {
		return nil, errors.WithStack(ErrHalideBuffer2DFloat32)
	}
	return (*C.halide_buffer_t)(buf), nil
}

func HalideFreeBuffer(buf *C.halide_buffer_t) {
	C.free_halide_buffer(buf)
}
