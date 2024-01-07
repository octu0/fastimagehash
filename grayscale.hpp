//go:build ignore

#include <Halide.h>
using namespace Halide;

Func grayscale(Func in, Expr width, Expr height);
std::tuple<Func, std::vector<Argument>> export_grayscale();
