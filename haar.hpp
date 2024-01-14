//go:build ignore

#include <Halide.h>
using namespace Halide;

Pipeline haar_x(Func input, Expr width, Expr height);
Pipeline haar_y(Func input, Expr width, Expr height);
Pipeline haar(Func input, Expr width, Expr height);
Func haar_hi(Func input, Expr width, Expr height);

std::tuple<Pipeline, std::vector<Argument>> export_haar_x();
std::tuple<Pipeline, std::vector<Argument>> export_haar_y();
std::tuple<Pipeline, std::vector<Argument>> export_haar();
std::tuple<Func, std::vector<Argument>> export_haar_hi();
