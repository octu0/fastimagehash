package fastimagehash

import (
	"sort"
)

func rgbaGrayscaleToGray(rgba []byte) []byte {
	out := make([]byte, 0, len(rgba)/4)
	for i := 0; i < len(rgba); i += 4 {
		out = append(out, rgba[i])
	}
	return out
}

func pickMedian(src []float32) float32 {
	tmp := make([]float32, len(src))
	copy(tmp, src)

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i] < tmp[j]
	})
	half := len(src) / 2
	return tmp[half]
}
