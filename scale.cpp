//go:build ignore

#include <Halide.h>
#include "scale.hpp"

using namespace Halide;

Func scale_kernel_box() {
  Var x("x");
  Func f = Func("scale_kernel_box");
  f(x) = select(abs(x) < 0.5f, 1.f, 0.f);

  f.compute_root();
  return f;
}

Func scale_kernel_linear() {
  Var x("x");
  Func f = Func("scale_kernel_linear");
  Expr xx = abs(x);
  f(x) = select(xx < 1.f, 1.f - xx, 0.f);

  f.compute_root();
  return f;
}

Func scale_kernel_gaussian() {
  Var x("x");
  Func f = Func("scale_kernel_gaussian");
  Expr xx = abs(x);
  Expr xx2 = fast_pow(0.5f, fast_pow(xx, 2.f));
  Expr base = fast_pow(0.5f, fast_pow(2, 2.f));
  f(x) = select(xx < 1.f, (xx2 - base) / (1 - base), 0.f);

  f.compute_root();
  return f;
}

Func scale_normal(
  Func input,
  Expr width, Expr height,
  Expr scale_width, Expr scale_height
) {
  Var x("x"), y("y"), ch("ch");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Region src_bounds = {{0, width},{0, height},{0, 4}};
  Func in = BoundaryConditions::constant_exterior(input, 0, src_bounds);

  Expr dx = cast<float>(width) / cast<float>(scale_width);
  Expr dy = cast<float>(height) / cast<float>(scale_height);

  Func fn = Func("scale_normal");
  Expr xx = cast<int>((x + 0.5f) * dx);
  Expr yy = cast<int>((y + 0.5f) * dy);

  fn(x, y, ch) = in(xx, yy, ch);

  fn.compute_at(in, ti)
    .store_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 8, 8)
    .fuse(xo, yo, ti)
    .parallel(ch)
    .parallel(ti, 8)
    .vectorize(xi, 8);
  return fn;
}

std::tuple<Func, std::vector<Argument>> export_scale_normal() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Param<int32_t> scale_width{"scale_width", 320};
  Param<int32_t> scale_height{"sale_height", 240};

  Func fn = scale_normal(src.in(), width, height, scale_width, scale_height);
  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height, scale_width, scale_height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

Func scale_by_kernel(
  Func input,
  Expr width, Expr height,
  Expr scale_width, Expr scale_height,
  Func kernel, Expr size,
  const char* name
) {
  Var x("x"), y("y"), ch("ch");
  Var s("s");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Region src_bounds = {{0, width},{0, height},{0, 4}};
  Func in = BoundaryConditions::constant_exterior(input, 0, src_bounds);

  Expr delta_w = cast<float>(width) / cast<float>(scale_width);
  Expr delta_h = cast<float>(height) / cast<float>(scale_height);
  Expr rate_w = max(1.0f, delta_w);
  Expr rate_h = max(1.0f, delta_h);
  Expr kernel_radius_w = rate_w * 1.0f;
  Expr kernel_radius_h = rate_h * 1.0f;
  RDom rd_scale = RDom(0, size, "rd_scale_box");

  Expr src_x = ((x + 0.5f) * delta_w) - 0.5f;
  Expr src_y = ((y + 0.5f) * delta_h) - 0.5f;
  Expr begin_x = cast<int>(ceil(src_x - kernel_radius_w));
  Expr begin_y = cast<int>(ceil(src_y - kernel_radius_h));
  begin_x = clamp(begin_x, 0, (width + 1) - size);
  begin_y = clamp(begin_y, 0, (height + 1) - size);

  Func kernel_val_x = Func("kernel_val_x"), kernel_val_y = Func("kernel_val_y");
  kernel_val_x(x, s) = kernel(cast<int>((s + begin_x - src_x) * rate_w));
  kernel_val_y(y, s) = kernel(cast<int>((s + begin_y - src_y) * rate_h));

  Func kernel_sum_x = Func("kernel_sum_x"), kernel_sum_y = Func("kernel_sum_y");
  kernel_sum_x(x) = sum(kernel_val_x(x, rd_scale));
  kernel_sum_y(y) = sum(kernel_val_y(y, rd_scale));

  Func kernel_x = Func("kernel_x"), kernel_y = Func("kernel_y");
  kernel_x(x, s) = kernel_val_x(x, s) / kernel_sum_x(x);
  kernel_y(y, s) = kernel_val_y(y, s) / kernel_sum_y(y);

  Func scale_y = Func("scale_y");
  Expr value = cast<float>(in(x, rd_scale + begin_y, ch));
  scale_y(x, y, ch) = sum(kernel_y(y, rd_scale) * value);

  Func scale_x = Func("scale_x");
  scale_x(x, y, ch) = sum(kernel_x(x, rd_scale) * scale_y(begin_x + rd_scale, y, ch));

  Func f = Func(name);
  Expr scaled = select(
    ch == 3, 255.0f,
    scale_x(x, y, ch)
  );
  f(x, y, ch) = cast<uint8_t>(scaled);

  kernel_val_x.compute_at(kernel_x, x)
    .vectorize(x);
  kernel_sum_x.compute_at(kernel_x, x)
    .vectorize(x);
  kernel_x.compute_root()
    .reorder(s, x)
    .vectorize(x, 8);

  kernel_val_y.compute_at(kernel_y, y)
    .vectorize(y, 8);
  kernel_sum_y.compute_at(kernel_y, y)
    .vectorize(y);
  kernel_y.compute_at(f, yi)
    .reorder(s, y)
    .vectorize(y, 8);

  f.compute_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 8, 8)
    .fuse(xo, yo, ti)
    .parallel(ch)
    .parallel(ti, 8)
    .vectorize(xi, 8);
  return f;
}

std::tuple<Func, std::vector<Argument>> export_scale_box() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Param<int32_t> scale_width{"scale_width", 320};
  Param<int32_t> scale_height{"sale_height", 240};

  Func fn = scale_by_kernel(src.in(), width, height, scale_width, scale_height, scale_kernel_box(), 1, "scale_box");
  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height, scale_width, scale_height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Func, std::vector<Argument>> export_scale_linear() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Param<int32_t> scale_width{"scale_width", 320};
  Param<int32_t> scale_height{"sale_height", 240};

  Func fn = scale_by_kernel(src.in(), width, height, scale_width, scale_height, scale_kernel_linear(), 1, "scale_linear");
  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height, scale_width, scale_height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Func, std::vector<Argument>> export_scale_gauss() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Param<int32_t> scale_width{"scale_width", 320};
  Param<int32_t> scale_height{"sale_height", 240};

  Func fn = scale_by_kernel(src.in(), width, height, scale_width, scale_height, scale_kernel_gaussian(), 1, "scale_gauss");
  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height, scale_width, scale_height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}
