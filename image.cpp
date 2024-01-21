//go:build ignore

#include <Halide.h>
#include "image.hpp"

using namespace Halide;

const Expr float_0 = cast<float>(0.f);
const Expr float_16 = cast<float>(16.f);
const Expr float_128 = cast<float>(128.f);
const Expr float_255 = cast<float>(255.f);

Func yuv_bt709_limited(Func yf, Func uf, Func vf, const char *name) {
  Var x("x"), y("y"), ch("ch");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Func f = Func(name);
  Expr yy = (yf(x, y) - float_16) * 1.164f;
  Expr uu = uf(x, y) - float_128;
  Expr vv = vf(x, y) - float_128;

  Expr r = yy + (1.793f * vv);
  Expr g = yy - (0.213f * uu) - (0.533f * vv);
  Expr b = yy + (2.112f * uu);

  Expr rr = clamp(r, float_0, float_255);
  Expr gg = clamp(g, float_0, float_255);
  Expr bb = clamp(b, float_0, float_255);

  Expr v = select(
    ch == 0, rr,       // R
    ch == 1, gg,       // G
    ch == 2, bb,       // B
    likely(float_255)  // A always 0xff
  );
  f(x, y, ch) = cast<uint8_t>(v);

  yf.compute_at(f, ch)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti, 16)
    .vectorize(xi, 32);

  uf.compute_at(f, ch)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti, 16)
    .vectorize(xi, 32);

  vf.compute_at(f, ch)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti, 16)
    .vectorize(xi, 32);
  return f;
}

Func yuv444_to_rgba_bt709(Func in_y, Func in_u, Func in_v) {
  Var x("x"), y("y");

  Func yf = Func("yuv444_in_y");
  Func uf = Func("yuv444_in_u");
  Func vf = Func("yuv444_in_v");
  yf(x, y) = cast<float>(in_y(x, y));
  uf(x, y) = cast<float>(in_u(x, y));
  vf(x, y) = cast<float>(in_v(x, y)); 
 
  return yuv_bt709_limited(yf, uf, vf, "yuv444_to_rgba");
}

Func yuv422_to_rgba_bt709(Func in_y, Func in_u, Func in_v) {
  Var x("x"), y("y");

  Func yf = Func("yuv422_in_y");
  Func uf = Func("yuv422_in_u");
  Func vf = Func("yuv422_in_v");
  yf(x, y) = cast<float>(in_y(x, y));
  uf(x, y) = cast<float>(in_u(x / 2, y));
  vf(x, y) = cast<float>(in_v(x / 2, y)); 

  return yuv_bt709_limited(yf, uf, vf, "yuv422_to_rgba");
}

Func yuv420_to_rgba_bt709(Func in_y, Func in_u, Func in_v) {
  Var x("x"), y("y");

  Func yf = Func("yuv420_in_y");
  Func uf = Func("yuv420_in_u");
  Func vf = Func("yuv420_in_v");
  yf(x, y) = cast<float>(in_y(x, y));
  uf(x, y) = cast<float>(in_u(x / 2, y / 2));
  vf(x, y) = cast<float>(in_v(x / 2, y / 2)); 

  return yuv_bt709_limited(yf, uf, vf, "yuv420_to_rgba");
}

std::tuple<Func, std::vector<Argument>> export_yuv444_to_rgba() {
  ImageParam y(UInt(8), 2, "y");
  ImageParam u(UInt(8), 2, "u");
  ImageParam v(UInt(8), 2, "v");

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};
  Param<int32_t> stride_y{"stride_y", 1920};
  Param<int32_t> stride_u{"stride_u", 1920};
  Param<int32_t> stride_v{"stride_v", 1920};

  y.dim(0).set_extent(width);
  y.dim(0).set_stride(1);
  y.dim(1).set_extent(height);
  y.dim(1).set_stride(stride_y);

  u.dim(0).set_extent(width);
  u.dim(0).set_stride(1);
  u.dim(1).set_extent(height);
  u.dim(1).set_stride(stride_u);

  v.dim(0).set_extent(width);
  v.dim(0).set_stride(1);
  v.dim(1).set_extent(height);
  v.dim(1).set_stride(stride_v);

  Func fn = yuv444_to_rgba_bt709(y.in(), u.in(), v.in());

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {y, u, v, stride_y, stride_u, stride_v, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Func, std::vector<Argument>> export_yuv422_to_rgba() {
  ImageParam y(UInt(8), 2, "y");
  ImageParam u(UInt(8), 2, "u");
  ImageParam v(UInt(8), 2, "v");

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};
  Param<int32_t> stride_y{"stride_y", 1920};
  Param<int32_t> stride_u{"stride_u", 960};
  Param<int32_t> stride_v{"stride_v", 960};

  y.dim(0).set_extent(width);
  y.dim(0).set_stride(1);
  y.dim(1).set_extent(height);
  y.dim(1).set_stride(stride_y);

  u.dim(0).set_extent(width / 2);
  u.dim(0).set_stride(1);
  u.dim(1).set_extent(height);
  u.dim(1).set_stride(stride_u);

  v.dim(0).set_extent(width / 2);
  v.dim(0).set_stride(1);
  v.dim(1).set_extent(height);
  v.dim(1).set_stride(stride_v);

  Func fn = yuv422_to_rgba_bt709(y.in(), u.in(), v.in());

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {y, u, v, stride_y, stride_u, stride_v, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Func, std::vector<Argument>> export_yuv420_to_rgba() {
  ImageParam y(UInt(8), 2, "y");
  ImageParam u(UInt(8), 2, "u");
  ImageParam v(UInt(8), 2, "v");

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};
  Param<int32_t> stride_y{"stride_y", 1920};
  Param<int32_t> stride_u{"stride_u", 960};
  Param<int32_t> stride_v{"stride_v", 960};

  y.dim(0).set_extent(width);
  y.dim(0).set_stride(1);
  y.dim(1).set_extent(height);
  y.dim(1).set_stride(stride_y);

  u.dim(0).set_extent(width / 2);
  u.dim(0).set_stride(1);
  u.dim(1).set_extent(height / 2);
  u.dim(1).set_stride(stride_u);

  v.dim(0).set_extent(width / 2);
  v.dim(0).set_stride(1);
  v.dim(1).set_extent(height / 2);
  v.dim(1).set_stride(stride_v);

  Func fn = yuv420_to_rgba_bt709(y.in(), u.in(), v.in());

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {y, u, v, stride_y, stride_u, stride_v, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}
