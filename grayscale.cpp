// +build ignore

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
  Expr a = cast<float>(in(x, y, 3));
  Expr sum = (r * gray_r_bt709) + (g * gray_g_bt709) + (b * gray_b_bt709);
  Expr sumi16 = cast<int16_t>(sum);
  Expr value = cast<uint8_t>(sumi16 >> 8);

  fn(x, y, ch) = cast<uint8_t>(255);
  fn(x, y, 0) = value; 
  fn(x, y, 1) = value; 
  fn(x, y, 2) = value; 
  fn(x, y, 3) = cast<uint8_t>(a); 

  // schedule
  fn.update(0).unscheduled();
  fn.update(1).unscheduled();
  fn.update(2).unscheduled();
  fn.update(3).unscheduled();
  fn.compute_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ch)
    .parallel(ti, 8)
    .vectorize(xi, 32);
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
