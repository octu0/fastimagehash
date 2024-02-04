# `fastimagehash`

[![MIT License](https://img.shields.io/github/license/octu0/fastimagehash)](https://github.com/octu0/fastimagehash/blob/master/LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/octu0/fastimagehash)](https://pkg.go.dev/github.com/octu0/fastimagehash)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/fastimagehash)](https://goreportcard.com/report/github.com/octu0/fastimagehash)
[![Releases](https://img.shields.io/github/v/release/octu0/fastimagehash)](https://github.com/octu0/fastimagehash/releases)

fast image hashing library for Go.

fastimagehash supports:

- Average Hashing
- Diffrence Hashing
- Perception Hashing
- Wavelet Hashing

## Installation

```bash
go get github.com/octu0/fastimagehash
```

## Benchmark

The basic score was based on [goimagehash](https://github.com/corona10/goimagehash).

```
goos: darwin
goarch: amd64
pkg: github.com/octu0/fastimagehash
cpu: Intel(R) Core(TM) i5-8210Y CPU @ 1.60GHz
BenchmarkAverageHash
BenchmarkAverageHash/goimagehash/8x8
BenchmarkAverageHash/goimagehash/8x8-4         	    2707	    424327 ns/op
BenchmarkAverageHash/goimagehash/32x32
BenchmarkAverageHash/goimagehash/32x32-4       	    1836	    570279 ns/op
BenchmarkAverageHash/fastimagehash/32x32
BenchmarkAverageHash/fastimagehash/32x32-4     	   39634	     29911 ns/op
BenchmarkDifferenceHash
BenchmarkDifferenceHash/goimagehash/9x8
BenchmarkDifferenceHash/goimagehash/9x8-4      	    2277	    456987 ns/op
BenchmarkDifferenceHash/goimagehash/33x32
BenchmarkDifferenceHash/goimagehash/33x32-4    	    1874	    592438 ns/op
BenchmarkDifferenceHash/fastimagehash/33x32
BenchmarkDifferenceHash/fastimagehash/33x32-4  	   33601	     34354 ns/op
BenchmarkPerceptionHash
BenchmarkPerceptionHash/goimagehash/64x64
BenchmarkPerceptionHash/goimagehash/64x64-4    	    1311	    813998 ns/op
BenchmarkPerceptionHash/fastimagehash/64x64
BenchmarkPerceptionHash/fastimagehash/64x64-4  	    7542	    146732 ns/op
BenchmarkWaveletHash
BenchmarkWaveletHash/fastimagehash/64x64
BenchmarkWaveletHash/fastimagehash/64x64-4     	    6045	    191712 ns/op
```

# License

MIT, see LICENSE file for details.
