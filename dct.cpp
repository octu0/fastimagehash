//go:build ignore

#include <Halide.h>
#include "dct.hpp"

using namespace Halide;

const Expr pi = cast<float>(3.14159265f);
const Expr isqrt2 = cast<float>(1.0f / sqrt(2.0f));

Func dct2d(Func input, Expr width, Expr height) {
  Var x("x"), y("y");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");

  Expr W = cast<float>(width);
  Expr H = cast<float>(height);

  Func cos_x = Func("cos_x");
  Expr x_a = sqrt(2.0f / W);
  Expr x_b = select(x == 0, isqrt2, 1.0f);
  Expr x_c = pi * (2 * x + 1) * y / (2 * W);
  Expr x_d = x_a * x_c * x_b;
  cos_x(x, y) = x_d;

  Func cos_y = Func("cos_y");
  Expr y_a = sqrt(2.0f / H);
  Expr y_b = select(y == 0, isqrt2, 1.0f);
  Expr y_c = pi * (2 * x + 1) * y / (2 * H);
  Expr y_d = y_a * y_c * y_b;
  cos_y(x, y) = y_d;

  Func fn = Func("dct2d");
  RDom rd_w = RDom(0, width);
  RDom rd_h = RDom(0, height);
  Expr dx = sum(input(x, rd_h) * cos_x(x, rd_h) * cos_y(x, rd_h));
  Expr dy = sum(input(rd_w, y) * cos_x(rd_w, y) * cos_y(rd_w, y));
  fn(x, y) =  dx * dy;

  fn.compute_at(input, y)
    .tile(x, y, xi, yi, 8, 8)
    .parallel(y)
    .vectorize(xi);
  return fn;
}

std::tuple<Func, std::vector<Argument>> export_dct2d() {
  ImageParam src(UInt(8), 2, "src");

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  src.dim(0).set_stride(1);
  src.dim(1).set_stride(width);

  Func fn = dct2d(src.in(), width, height);
  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(1);
  out.dim(1).set_stride(width);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}
