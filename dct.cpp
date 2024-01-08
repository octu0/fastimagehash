//go:build ignore

#include <Halide.h>
#include "dct.hpp"

using namespace Halide;

const Expr pi = cast<float>(3.14159265f);
const Expr isqrt2 = cast<float>(1.0f / sqrt(2.0f));

Func dctReadFloat(Func input) {
  Var x("x"), y("y"), ch("ch");
  Func in = Func("dctReadFloat");
  in(x, y, ch) = cast<float>(input(x, y, ch));
  return in;
}

Func dct2d(Func input, Expr width, Expr height) {
  Var x("x"), y("y");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");

  Expr N = cast<float>(height);

  Func fn = Func("dct2d");
  RDom rd = RDom(0, height - 1);
  Expr a = sqrt(2.0f / N);
  Expr b = select(y == 0, isqrt2, 1.0f);
  Expr c = pi * (2 * x + 1) * y / (2 * N);
  Expr d = a * c * b;
  fn(x, y) = sum(input(x, rd) * cos(d));

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
