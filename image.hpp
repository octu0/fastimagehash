//go:build ignore

#include <Halide.h>
using namespace Halide;

Func yuv444_to_rgba_bt709(Func in_y, Func in_u, Func in_v);
Func yuv422_to_rgba_bt709(Func in_y, Func in_u, Func in_v);
Func yuv420_to_rgba_bt709(Func in_y, Func in_u, Func in_v);

std::tuple<Func, std::vector<Argument>> export_yuv444_to_rgba();
std::tuple<Func, std::vector<Argument>> export_yuv422_to_rgba();
std::tuple<Func, std::vector<Argument>> export_yuv420_to_rgba();
