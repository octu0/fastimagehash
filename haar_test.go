package fastimagehash

import (
	"testing"
)

func TestHaarX(t *testing.T) {
	signal := []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5,
		3, 2, 2, 2, 3, 3, 3, 3, 2, 2, 2, 2, 1, 1, 1, 1,
		4, 2, 2, 2, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3,
		5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5,
		6, 1, 2, 3, 4, 5, 6, 7, 6, 5, 4, 3, 4, 5, 6, 7,
		7, 0, 7, 0, 7, 0, 7, 0, 7, 0, 7, 0, 7, 0, 7, 0,
		8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
		9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	lo, hi, err := haarX(signal, 16, 16)
	if err != nil {
		t.Errorf("no error: %+v", err)
	}

	btof32 := func(data []byte) []float32 {
		f := make([]float32, len(data))
		for i := 0; i < len(data); i += 1 {
			f[i] = float32(data[i])
		}
		return f
	}
	expectLo := make([][]float32, 0, 16)
	expectHi := make([][]float32, 0, 16)
	for i := 0; i < len(signal); i += 16 {
		sig := btof32(signal[i : i+16])
		eLo, eHi := nativeHaar1D(sig)
		expectLo = append(expectLo, eLo)
		expectHi = append(expectHi, eHi)
	}
	equal := func(a, b float32) bool {
		r := a - b
		if r == 0 {
			return true
		}
		// margin of error
		if 0 < r && 0.000001999 < r {
			return false
		}
		if r < 0 && r < -0.000001999 {
			return false
		}
		return true
	}
	equalSlice := func(a, b []float32) bool {
		if len(a) != len(b) {
			return false
		}
		for i := 0; i < len(a); i += 1 {
			if equal(a[i], b[i]) != true {
				t.Logf("a %v <> b %v: %3.8f", a[i], b[i], a[i]-b[i])
				return false
			}
		}
		return true
	}

	if equalSlice(expectLo[0], lo[0:8]) != true {
		t.Errorf("expect %v <> actual %v", expectLo[0], lo[0:8])
	}
	if equalSlice(expectHi[0], hi[0:8]) != true {
		t.Errorf("expect %v <> actual %v", expectHi[0], hi[0:8])
	}

	if equalSlice(expectLo[1], lo[8:16]) != true {
		t.Errorf("expect %v <> actual %v", expectLo[1], lo[8:16])
	}
	if equalSlice(expectHi[1], hi[8:16]) != true {
		t.Errorf("expect %v <> actual %v", expectHi[1], hi[8:16])
	}

	if equalSlice(expectLo[2], lo[16:24]) != true {
		t.Errorf("expect %v <> actual %v", expectLo[2], lo[16:24])
	}
	if equalSlice(expectHi[2], hi[16:24]) != true {
		t.Errorf("expect %v <> actual %v", expectHi[2], hi[16:24])
	}

	if equalSlice(expectLo[3], lo[24:32]) != true {
		t.Errorf("expect %v <> actual %v", expectLo[3], lo[24:32])
	}
	if equalSlice(expectHi[3], hi[24:32]) != true {
		t.Errorf("expect %v <> actual %v", expectHi[3], hi[24:32])
	}

	if equalSlice(expectLo[4], lo[32:40]) != true {
		t.Errorf("expect %v <> actual %v", expectLo[4], lo[32:40])
	}
	if equalSlice(expectHi[4], hi[32:40]) != true {
		t.Errorf("expect %v <> actual %v", expectHi[4], hi[32:40])
	}

	if equalSlice(expectLo[5], lo[40:48]) != true {
		t.Errorf("expect %v <> actual %v", expectLo[5], lo[40:48])
	}
	if equalSlice(expectHi[5], hi[40:48]) != true {
		t.Errorf("expect %v <> actual %v", expectHi[5], hi[40:48])
	}

	if equalSlice(expectLo[6], lo[48:56]) != true {
		t.Errorf("expect %v <> actual %v", expectLo[6], lo[48:56])
	}
	if equalSlice(expectHi[6], hi[48:56]) != true {
		t.Errorf("expect %v <> actual %v", expectHi[6], hi[48:56])
	}

	p := 7
	for i := 56; i < 128; i += 8 {
		if equalSlice(expectLo[p], lo[i:i+8]) != true {
			t.Errorf("[%d] expect %v <> actual %v", p, expectLo[p], lo[i:i+8])
		}
		if equalSlice(expectHi[p], hi[i:i+8]) != true {
			t.Errorf("[%d] expect %v <> actual %v", p, expectHi[p], hi[i:i+8])
		}
		p += 1
	}
}

func TestHaarY(t *testing.T) {
	signal := []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		2, 2, 2, 2, 0, 1, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 2, 2, 2, 0, 2, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		4, 2, 2, 2, 0, 3, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		5, 3, 3, 1, 0, 4, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		6, 3, 3, 1, 0, 5, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		7, 3, 3, 1, 0, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		8, 3, 3, 1, 0, 7, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		9, 4, 2, 2, 0, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		10, 4, 2, 2, 0, 5, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		11, 4, 2, 2, 0, 4, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		12, 4, 2, 2, 0, 3, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		13, 5, 1, 3, 0, 4, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		14, 5, 1, 3, 0, 5, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		15, 5, 1, 3, 0, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		16, 5, 1, 3, 5, 7, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	lo, hi, err := haarY(signal, 16, 16)
	if err != nil {
		t.Errorf("no error: %+v", err)
	}

	btof32 := func(data []byte) []float32 {
		f := make([]float32, len(data))
		for i := 0; i < len(data); i += 1 {
			f[i] = float32(data[i])
		}
		return f
	}
	equal := func(a, b float32) bool {
		r := a - b
		if r == 0 {
			return true
		}
		// margin of error
		if 0 < r && 0.000001999 < r {
			return false
		}
		if r < 0 && r < -0.000001999 {
			return false
		}
		return true
	}
	equalSlice := func(a, b []float32) bool {
		if len(a) != len(b) {
			return false
		}
		for i := 0; i < len(a); i += 1 {
			if equal(a[i], b[i]) != true {
				t.Logf("a %v <> b %v: %3.8f", a[i], b[i], a[i]-b[i])
				return false
			}
		}
		return true
	}

	for i := 0; i < 16; i += 1 {
		aLo := []float32{
			lo[i+0],
			lo[i+16],
			lo[i+32],
			lo[i+48],
			lo[i+64],
			lo[i+80],
			lo[i+96],
			lo[i+112],
		}
		aHi := []float32{
			hi[i+0],
			hi[i+16],
			hi[i+32],
			hi[i+48],
			hi[i+64],
			hi[i+80],
			hi[i+96],
			hi[i+112],
		}
		inY := []byte{
			signal[i+0],
			signal[i+16],
			signal[i+32],
			signal[i+48],
			signal[i+64],
			signal[i+80],
			signal[i+96],
			signal[i+112],
			signal[i+128],
			signal[i+144],
			signal[i+160],
			signal[i+176],
			signal[i+192],
			signal[i+208],
			signal[i+224],
			signal[i+240],
		}
		expectLo, expectHi := nativeHaar1D(btof32(inY))
		if equalSlice(expectLo, aLo) != true {
			t.Errorf("[%d] expect %v <> actual %v", i, expectLo, aLo)
		}
		if equalSlice(expectHi, aHi) != true {
			t.Errorf("[%d] expect %v <> actual %v", i, expectHi, aHi)
		}
	}
}

func TestHaarXY(t *testing.T) {
	signal := []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		2, 2, 2, 2, 0, 1, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 2, 2, 2, 0, 2, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		4, 2, 2, 2, 0, 3, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		5, 3, 3, 1, 0, 4, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		6, 3, 3, 1, 0, 5, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		7, 3, 3, 1, 0, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		8, 3, 3, 1, 0, 7, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		9, 4, 2, 2, 0, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		10, 4, 2, 2, 0, 5, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		11, 4, 2, 2, 0, 4, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		12, 4, 2, 2, 0, 3, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		13, 5, 1, 3, 0, 4, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		14, 5, 1, 3, 0, 5, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		15, 5, 1, 3, 0, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0,
		16, 5, 1, 3, 5, 7, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	lo, hi, err := haar(signal, 16, 16)
	if err != nil {
		t.Errorf("no error: %+v", err)
	}

	btof32 := func(data []byte) []float32 {
		f := make([]float32, len(data))
		for i := 0; i < len(data); i += 1 {
			f[i] = float32(data[i])
		}
		return f
	}
	equal := func(a, b float32) bool {
		r := a - b
		if r == 0 {
			return true
		}
		// margin of error
		if 0 < r && 0.000001999 < r {
			return false
		}
		if r < 0 && r < -0.000001999 {
			return false
		}
		return true
	}
	equalSlice := func(a, b []float32) bool {
		if len(a) != len(b) {
			return false
		}
		for i := 0; i < len(a); i += 1 {
			if equal(a[i], b[i]) != true {
				t.Logf("a %v <> b %v: %3.8f", a[i], b[i], a[i]-b[i])
				return false
			}
		}
		return true
	}

	eLoXY, eHiXY := nativeHaar2D([][]float32{
		btof32(signal[0:16]),
		btof32(signal[16:32]),
		btof32(signal[32:48]),
		btof32(signal[48:64]),
		btof32(signal[64:80]),
		btof32(signal[80:96]),
		btof32(signal[96:112]),
		btof32(signal[112:128]),
		btof32(signal[128:144]),
		btof32(signal[144:160]),
		btof32(signal[160:176]),
		btof32(signal[176:192]),
		btof32(signal[192:208]),
		btof32(signal[208:224]),
		btof32(signal[224:240]),
		btof32(signal[240:256]),
	})
	p := 0
	for i := 0; i < 64; i += 8 {
		aLo := lo[i : i+8]
		aHi := hi[i : i+8]

		if equalSlice(eLoXY[p], aLo) != true {
			t.Logf("[%d] expect %v <> actual %v", p, eLoXY[p], aLo)
		}
		if equalSlice(eHiXY[p], aHi) != true {
			t.Logf("[%d] expect %v <> actual %v", p, eHiXY[p], aHi)
		}
		p += 1
	}
}
