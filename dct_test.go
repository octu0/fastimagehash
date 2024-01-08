package fastimagehash

import (
	"testing"
)

func TestDCT2d(t *testing.T) {
	in := []byte{
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16,
		20, 21, 22, 23, 24, 25, 26, 27,
		30, 31, 32, 33, 34, 35, 36, 37,
	}
	out, err := dct2d(in, 8, 8)
	if err != nil {
		t.Errorf("no error %+v", err)
	}
	t.Logf("%+v", out)
}
