package fastimagehash

import (
	"testing"
)

func BenchmarkWaveletHash(b *testing.B) {
	b.Run("fastimagehash/64x64", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			hash, err := WaveletHash(rgba)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = hash.Hex()
		}
	})
}

func TestBaz(t *testing.T) {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}
	ahash, err := WaveletHash(rgba)
	if err != nil {
		panic(err)
	}
	println(ahash.Hex())

	{
		rgba, err := pngToRGBA(testPngBlendImg)
		if err != nil {
			panic(err)
		}
		a, err := WaveletHash(rgba)
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
		a, err := WaveletHash(rgba)
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
		a, err := WaveletHash(rgba)
		if err != nil {
			panic(err)
		}
		println("cat = ", a.Hex())
		println("distance3=", ahash.Distance(a))
	}
}
