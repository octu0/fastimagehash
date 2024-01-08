package fastimagehash

import (
	"testing"

	"github.com/pkg/errors"
)

func ExampleHash1024() {
	h := Hash1024([16]uint64{
		1, 8, 16, 32, 64, 128, 256, 512,
		1, 2, 3, 4, 5, 6, 7,
	})
	println(h.Hex())

	// Output:
}

func TestHash1024FromString(t *testing.T) {
	t.Run("err/size<256", func(tt *testing.T) {
		_, err := Hash1024FromString("deadbeaf")
		if errors.Is(err, ErrHash1024InvalidSize) != true {
			tt.Errorf("invalid size error")
		}
	})
	t.Run("restore/hex", func(tt *testing.T) {
		hex := "01000000bc01b000e27de200ebfe4000cfffc000effe2806e7f64f0097e61f809f663f80ff863fc1f787bf80ee9600c23eb001810346001981ccc11d80200017e1a401af11c210cf81e1364707d0203e61b80416f94080a6fb50003f1110100393681001f1640a8af5c502c9fd208a89fc01a081fc460101e45400127ef08260"
		hash, err := Hash1024FromString(hex)
		if err != nil {
			tt.Errorf("no error")
		}
		if hash.Hex() != hex {
			tt.Errorf("restore failed expect:%s != actual:%s", hex, hash.Hex())
		}
	})
}
