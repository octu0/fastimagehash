//go:build ignore

#include <Halide.h>
using namespace Halide;

Func dct2d(Func input, Expr width, Expr height);
std::tuple<Func, std::vector<Argument>> export_dct2d();
