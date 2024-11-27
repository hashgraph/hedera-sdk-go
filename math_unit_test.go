//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"math/big"
	"testing"
)

func TestPaddedBigBytes(t *testing.T) {
	tests := []struct {
		num    *big.Int
		n      int
		result []byte
	}{
		{num: big.NewInt(0), n: 4, result: []byte{0, 0, 0, 0}},
		{num: big.NewInt(1), n: 4, result: []byte{0, 0, 0, 1}},
		{num: big.NewInt(512), n: 4, result: []byte{0, 0, 2, 0}},
	}
	for _, test := range tests {
		if result := ToPaddedBytes(test.num, test.n); !bytes.Equal(result, test.result) {
			t.Errorf("PaddedBigBytes(%d, %d) = %v, want %v", test.num, test.n, result, test.result)
		}
	}
}

func TestS256(t *testing.T) {
	tests := []struct{ x, y *big.Int }{
		{x: big.NewInt(0), y: big.NewInt(0)},
		{x: big.NewInt(1), y: big.NewInt(1)},
		{x: big.NewInt(2), y: big.NewInt(2)},
		{
			x: new(big.Int).Sub(PowerOfBig(2, 255), big.NewInt(1)),
			y: new(big.Int).Sub(PowerOfBig(2, 255), big.NewInt(1)),
		},
		{
			x: PowerOfBig(2, 255),
			y: new(big.Int).Neg(PowerOfBig(2, 255)),
		},
		{
			x: new(big.Int).Sub(PowerOfBig(2, 256), big.NewInt(1)),
			y: big.NewInt(-1),
		},
		{
			x: new(big.Int).Sub(PowerOfBig(2, 256), big.NewInt(2)),
			y: big.NewInt(-2),
		},
	}
	for _, test := range tests {
		if y := ToSigned256(test.x); y.Cmp(test.y) != 0 {
			t.Errorf("S256(%x) = %x, want %x", test.x, y, test.y)
		}
	}
}

func TestU256Bytes(t *testing.T) {
	ubytes := make([]byte, 32)
	ubytes[31] = 1

	unsigned := To256BitBytes(big.NewInt(1))
	if !bytes.Equal(unsigned, ubytes) {
		t.Errorf("expected %x got %x", ubytes, unsigned)
	}
}

func TestU256(t *testing.T) {
	tests := []struct{ x, y *big.Int }{
		{x: big.NewInt(0), y: big.NewInt(0)},
		{x: big.NewInt(1), y: big.NewInt(1)},
		{x: PowerOfBig(2, 255), y: PowerOfBig(2, 255)},
		{x: PowerOfBig(2, 256), y: big.NewInt(0)},
		{x: new(big.Int).Add(PowerOfBig(2, 256), big.NewInt(1)), y: big.NewInt(1)},
		// negative values
		{x: big.NewInt(-1), y: new(big.Int).Sub(PowerOfBig(2, 256), big.NewInt(1))},
		{x: big.NewInt(-2), y: new(big.Int).Sub(PowerOfBig(2, 256), big.NewInt(2))},
		{x: PowerOfBig(2, -255), y: big.NewInt(1)},
	}
	for _, test := range tests {
		if y := To256Bit(new(big.Int).Set(test.x)); y.Cmp(test.y) != 0 {
			t.Errorf("U256(%x) = %x, want %x", test.x, y, test.y)
		}
	}
}
