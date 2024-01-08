package fastimagehash

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/bits"

	"github.com/pkg/errors"
)

var (
	ErrHash1024InvalidSize = errors.New("invalid hash1024 size")
)

type Hash1024 [16]uint64

func (h Hash1024) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	for i := 0; i < 16; i += 1 {
		binary.Write(buf, binary.BigEndian, h[i])
	}
	return buf.Bytes()
}

func (h Hash1024) Hex() string {
	return hex.EncodeToString(h.Bytes())
}

func (h Hash1024) String() string {
	return h.Hex()
}

func (h Hash1024) Distance(o Hash1024) int {
	distance := 0
	for i := 0; i < 16; i += 1 {
		lh := h[i]
		rh := o[i]
		hamming := uint64(lh ^ rh)
		distance += bits.OnesCount64(hamming)
	}
	return distance
}

func Hash1024FromBytes(data []byte) (Hash1024, error) {
	if len(data) != 128 {
		return Hash1024{}, errors.Wrapf(ErrHash1024InvalidSize, "len(bytes) != 128")
	}
	hash := Hash1024{}
	for i := 0; i < 16; i += 1 {
		hash[i] = binary.BigEndian.Uint64(data[:8])
		data = data[8:]
	}
	return hash, nil
}

func Hash1024FromString(s string) (Hash1024, error) {
	if len(s) != 256 {
		return Hash1024{}, errors.Wrapf(ErrHash1024InvalidSize, "len(hex) != 256")
	}
	data, err := hex.DecodeString(s)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}
	hash, err := Hash1024FromBytes(data)
	if err != nil {
		return Hash1024{}, errors.WithStack(err)
	}
	return hash, nil
}
