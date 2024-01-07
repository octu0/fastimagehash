//go:build ignore

#include <Halide.h>
using namespace Halide;

Func scale_normal(Func in, Expr width, Expr height, Expr scale_width, Expr scale_height);
Func scale_by_kernel(
  Func in,
  Expr width, Expr height,
  Expr scale_width, Expr scale_height,
  Func kernel, Expr size,
  const char *name
);

std::tuple<Func, std::vector<Argument>> export_scale_normal();
std::tuple<Func, std::vector<Argument>> export_scale_box();
std::tuple<Func, std::vector<Argument>> export_scale_linear();
std::tuple<Func, std::vector<Argument>> export_scale_gauss();
