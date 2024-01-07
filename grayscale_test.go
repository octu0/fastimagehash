package fastimagehash

import (
	"testing"
)

// https://github.com/octu0/go-intrin#benchmark
func BenchmarkGrayscale(b *testing.B) {
	b.Run("halide", func(tb *testing.B) {
		img, err := pngToRGBA(testPngImg)
		if err != nil {
			tb.Fatalf("%+v", err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			_, _ = grayscale(img)
		}
	})
}

func ExampleGrayscale() {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}

	gray, err := grayscale(rgba)
	if err != nil {
		panic(err)
	}

	path, err := saveImage(gray)
	if err != nil {
		panic(err)
	}
	println("grayscale = ", path)

	// Output:
}
