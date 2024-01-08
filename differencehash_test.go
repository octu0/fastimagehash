package fastimagehash

import (
	"testing"

	"github.com/corona10/goimagehash"
)

func BenchmarkDifferenceHash(b *testing.B) {
	b.Run("goimagehash/9x8", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			hash, err := goimagehash.DifferenceHash(rgba)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = hash.ToString()
		}
	})
	b.Run("goimagehash/33x32", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			hash, err := goimagehash.ExtDifferenceHash(rgba, 33, 32)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = hash.ToString()
		}
	})
	b.Run("fastimagehash/33x32", func(tb *testing.B) {
		rgba, err := pngToRGBA(testPngImg)
		if err != nil {
			panic(err)
		}
		tb.ResetTimer()
		for i := 0; i < tb.N; i += 1 {
			hash, err := DifferenceHash(rgba)
			if err != nil {
				tb.Fatalf("%+v", err)
			}
			_ = hash.Hex()
		}
	})
}

func TestFoo(t *testing.T) {
	rgba, err := pngToRGBA(testPngImg)
	if err != nil {
		panic(err)
	}
	dhash, err := DifferenceHash(rgba)
	if err != nil {
		panic(err)
	}
	println(dhash.Hex())

	{
		rgba, err := pngToRGBA(testPngBlendImg)
		if err != nil {
			panic(err)
		}
		d, err := DifferenceHash(rgba)
		if err != nil {
			panic(err)
		}
		println("blend = ", d.Hex())
		println("distance1=", dhash.Distance(d))
	}
	{
		rgba, err := pngToRGBA(testPngBlurImg)
		if err != nil {
			panic(err)
		}
		d, err := DifferenceHash(rgba)
		if err != nil {
			panic(err)
		}
		println("blur = ", d.Hex())
		println("distance2=", dhash.Distance(d))
	}
	{
		rgba, err := pngToRGBA(testPngCatA1Img)
		if err != nil {
			panic(err)
		}
		d, err := DifferenceHash(rgba)
		if err != nil {
			panic(err)
		}
		println("cat = ", d.Hex())
		println("distance3=", dhash.Distance(d))
	}

	edhash, err := goimagehash.ExtDifferenceHash(rgba, 33, 32)
	if err != nil {
		panic(err)
	}
	println(edhash.ToString())
}
