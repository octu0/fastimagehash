//go:build ignore

#include <Halide.h>
#include "haar.hpp"

using namespace Halide;

const Expr sqrt2d = cast<double>(1.4142135623730951f); // sqrt(2)
const Expr sqrt2 = cast<float>(1.4142135f); // sqrt(2)

Func read_xy_float(Func input, const char *name) {
  Var x("x"), y("y");

  Func in = Func(name);
  in(x, y) = cast<float>(input(x, y));
  return in;
}

Pipeline haar_x(Func input, Expr width, Expr height) {
  Var x("x"), y("y");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Region src_bounds = {{0, width},{0, height}};
  Func rf = read_xy_float(input, "haar_x_in");
  Func in = BoundaryConditions::constant_exterior(rf, 0.f, src_bounds);

  Func haar_x_lo = Func("haar_x_lo");
  haar_x_lo(x, y) = (in(2 * x, y) + in(2 * x + 1, y)) / sqrt2;
  Func haar_x_hi = Func("haar_x_hi");
  haar_x_hi(x, y) = (in(2 * x, y) - in(2 * x + 1, y)) / sqrt2;

  in.compute_root();
  haar_x_lo.compute_at(in, ti)
    .store_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti)
    .vectorize(xi, 32);
  haar_x_hi.compute_at(in, ti)
    .store_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti)
    .vectorize(xi, 32);

  return Pipeline({ haar_x_lo, haar_x_hi });
}

Pipeline haar_y(Func input, Expr width, Expr height) {
  Var x("x"), y("y");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Region src_bounds = {{0, width},{0, height}};
  Func rf = read_xy_float(input, "haar_y_in");
  Func in = BoundaryConditions::constant_exterior(rf, 0.f, src_bounds);

  Func haar_y_lo = Func("haar_y_lo");
  haar_y_lo(x, y) = (in(x, 2 * y) + in(x, 2 * y + 1)) / sqrt2;
  Func haar_y_hi = Func("haar_y_hi");
  haar_y_hi(x, y) = (in(x, y * 2) - in(x, 2 * y + 1)) / sqrt2;

  in.compute_root();

  haar_y_lo.compute_at(in, ti)
    .store_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti)
    .vectorize(xi, 32);
  haar_y_hi.compute_at(in, ti)
    .store_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti)
    .vectorize(xi, 32);

  return Pipeline({ haar_y_lo, haar_y_hi });
}

Pipeline haar(Func input, Expr width, Expr height) {
  Var x("x"), y("y");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Region src_bounds = {{0, width},{0, height}};
  Func rf = read_xy_float(input, "haar_xy_in");
  Func in = BoundaryConditions::constant_exterior(rf, 0.f, src_bounds);

  Func haar_x_lo = Func("haar_x_lo");
  haar_x_lo(x, y) = (in(2 * x, y) + in(2 * x + 1, y)) / sqrt2;
  Func haar_x_hi = Func("haar_x_hi");
  haar_x_hi(x, y) = (in(2 * x, y) - in(2 * x + 1, y)) / sqrt2;

  Func haar_xy_lo = Func("haar_xy_lo");
  haar_xy_lo(x, y) = (haar_x_lo(x, 2 * y) + haar_x_lo(x, 2 * y + 1)) / sqrt2;
  Func haar_xy_hi = Func("haar_xy_hi");
  haar_xy_hi(x, y) = (haar_x_hi(x, 2 * y) - haar_x_hi(x, 2 * y + 1)) / sqrt2;

  in.compute_root();

  haar_xy_lo.compute_at(haar_x_lo, ti)
    .store_at(haar_x_lo, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti)
    .vectorize(xi, 32);
  haar_xy_hi.compute_at(haar_x_hi, ti)
    .store_at(haar_x_hi, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti)
    .vectorize(xi, 32);

  return Pipeline({ haar_xy_lo, haar_xy_hi });
}

Func haar_hi(Func input, Expr width, Expr height) {
  Var x("x"), y("y");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Region src_bounds = {{0, width},{0, height}};
  Func rf = read_xy_float(input, "haar_hi_in");
  Func in = BoundaryConditions::constant_exterior(rf, 0.f, src_bounds);

  Func haar_x_hi = Func("haar_x_hi");
  haar_x_hi(x, y) = (in(2 * x, y) - in(2 * x + 1, y)) / sqrt2;

  Func haar_xy_hi = Func("haar_xy_hi");
  haar_xy_hi(x, y) = (haar_x_hi(x, 2 * y) - haar_x_hi(x, 2 * y + 1)) / sqrt2;

  in.compute_root();

  haar_xy_hi.compute_at(haar_x_hi, ti)
    .store_at(haar_x_hi, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ti)
    .vectorize(xi, 32);

  return haar_xy_hi;
}

std::tuple<Pipeline, std::vector<Argument>> export_haar_x() {
  ImageParam src(UInt(8), 2, "src");

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  src.dim(0).set_stride(1);
  src.dim(1).set_stride(width);

  Pipeline fn = haar_x(src.in(), width, height);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Pipeline, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Pipeline, std::vector<Argument>> export_haar_y() {
  ImageParam src(UInt(8), 2, "src");

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  src.dim(0).set_stride(1);
  src.dim(1).set_stride(width);

  Pipeline fn = haar_y(src.in(), width, height);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Pipeline, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Pipeline, std::vector<Argument>> export_haar() {
  ImageParam src(UInt(8), 2, "src");

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  src.dim(0).set_stride(1);
  src.dim(1).set_stride(width);

  Pipeline fn = haar(src.in(), width, height);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Pipeline, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Func, std::vector<Argument>> export_haar_hi() {
  ImageParam src(UInt(8), 2, "src");

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  src.dim(0).set_stride(1);
  src.dim(1).set_stride(width);

  Func fn = haar_hi(src.in(), width, height);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}
