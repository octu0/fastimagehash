package fastimagehash

import (
	"testing"

	"github.com/corona10/goimagehash"
)

func BenchmarkAverageHash(b *testing.B) {
	b.Run("goimagehash/8x8", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			ahash, err := goimagehash.AverageHash(rgba)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = ahash.ToString()
		}
	})
	b.Run("goimagehash/32x32", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			ahash, err := goimagehash.ExtAverageHash(rgba, 32, 32)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = ahash.ToString()
		}
	})
	b.Run("fastimagehash/32x32", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			ahash, err := AverageHash(rgba)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = ahash.Hex()
		}
	})
}

func TestHoge(t *testing.T) {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}
	ahash, err := AverageHash(rgba)
	if err != nil {
		panic(err)
	}
	println(ahash.Hex())

	{
		rgba, err := pngToRGBA(testPngBlendImg)
		if err != nil {
			panic(err)
		}
		a, err := AverageHash(rgba)
		if err != nil {
			panic(err)
		}
		println("blend = ", a.Hex())
		println("distance1=", ahash.Distance(a))
	}
	{
		rgba, err := pngToRGBA(testPngBlurImg)
		if err != nil {
			panic(err)
		}
		a, err := AverageHash(rgba)
		if err != nil {
			panic(err)
		}
		println("blur = ", a.Hex())
		println("distance2=", ahash.Distance(a))
	}
	{
		rgba, err := pngToRGBA(testPngCatA1Img)
		if err != nil {
			panic(err)
		}
		a, err := AverageHash(rgba)
		if err != nil {
			panic(err)
		}
		println("cat = ", a.Hex())
		println("distance3=", ahash.Distance(a))
	}

	eahash, err := goimagehash.ExtAverageHash(rgba, 32, 32)
	if err != nil {
		panic(err)
	}
	println(eahash.ToString())
}
