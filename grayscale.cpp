//go:build ignore

#include <Halide.h>
#include "grayscale.hpp"

using namespace Halide;

const Expr gray_r_bt709 = cast<float>(0.2126f);
const Expr gray_g_bt709 = cast<float>(0.7152f);
const Expr gray_b_bt709 = cast<float>(0.0722f);

Func grayscale(Func input, Expr width, Expr height) {
  Var x("x"), y("y"), ch("ch");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Region src_bounds = {{0, width},{0, height},{0, 4}};
  Func in = BoundaryConditions::repeat_edge(input, src_bounds);

  Func fn = Func("grayscale");
  Expr r = cast<float>(in(x, y, 0));
  Expr g = cast<float>(in(x, y, 1));
  Expr b = cast<float>(in(x, y, 2));
  Expr a = in(x, y, 3);
  Expr sum = (r * gray_r_bt709) + (g * gray_g_bt709) + (b * gray_b_bt709);
  Expr value = cast<uint8_t>(sum);

  fn(x, y, ch) = select(
    ch == 3, cast<uint8_t>(255),
    value
  );
  // schedule
  fn.compute_at(in, ti)
    .store_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 8, 8)
    .fuse(xo, yo, ti)
    .parallel(ch)
    .parallel(ti, 8)
    .vectorize(xi);
  return fn;
}

std::tuple<Func, std::vector<Argument>> export_grayscale() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Func fn = grayscale(src.in(), width, height);

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}
