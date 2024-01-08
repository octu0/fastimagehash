package fastimagehash

import (
	"testing"

	"github.com/corona10/goimagehash"
)

func BenchmarkPerceptionHash(b *testing.B) {
	b.Run("goimagehash/64x64", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			hash, err := goimagehash.PerceptionHash(rgba)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = hash.ToString()
		}
	})
	b.Run("fastimagehash/64x64", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			hash, err := PerceptionHash(rgba)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = hash.Hex()
		}
	})
}

func TestBar(t *testing.T) {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}
	hash, err := PerceptionHash(rgba)
	if err != nil {
		panic(err)
	}
	println(hash.Hex())

	{
		rgba, err := pngToRGBA(testPngBlendImg)
		if err != nil {
			panic(err)
		}
		h, err := PerceptionHash(rgba)
		if err != nil {
			panic(err)
		}
		println("blend = ", h.Hex())
		println("distance1=", hash.Distance(h))
	}
	{
		rgba, err := pngToRGBA(testPngBlurImg)
		if err != nil {
			panic(err)
		}
		h, err := PerceptionHash(rgba)
		if err != nil {
			panic(err)
		}
		println("blur = ", h.Hex())
		println("distance2=", hash.Distance(h))
	}
	{
		rgba, err := pngToRGBA(testPngCatA1Img)
		if err != nil {
			panic(err)
		}
		h, err := PerceptionHash(rgba)
		if err != nil {
			panic(err)
		}
		println("cat = ", h.Hex())
		println("distance3=", hash.Distance(h))
	}

	phash, err := goimagehash.PerceptionHash(rgba)
	if err != nil {
		panic(err)
	}
	println(phash.ToString())
}
