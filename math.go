package hiero

// SPDX-License-Identifier: Apache-2.0

import "math/big"

const (
	// BitsPerWord represents the number of bits in a big.Word
	BitsPerWord = 32 << (uint64(^big.Word(0)) >> 63)
	// BytesPerWord represents the number of bytes in a big.Word
	BytesPerWord = BitsPerWord / 8
)

// Various big integer constants for limits.
var (
	bigLimit255       = PowerOfBig(2, 255)
	bigLimit256       = PowerOfBig(2, 256)
	bigLimit256Minus1 = new(big.Int).Sub(bigLimit256, big.NewInt(1))
	bigLimit63        = PowerOfBig(2, 63)
	MaxBig256Value    = new(big.Int).Set(bigLimit256Minus1)
	MaxBig63Value     = new(big.Int).Sub(bigLimit63, big.NewInt(1))
)

// PowerOfBig computes a^b as a big integer.
func PowerOfBig(a, b int64) *big.Int {
	result := big.NewInt(a)
	return result.Exp(result, big.NewInt(b), nil)
}

// To256Bit truncates the input number to a 256-bit two's complement. This modifies the input.
func To256Bit(x *big.Int) *big.Int {
	return x.And(x, bigLimit256Minus1)
}

// To256BitBytes converts a big integer into a 32-byte slice in big-endian order.
func To256BitBytes(n *big.Int) []byte {
	return ToPaddedBytes(To256Bit(n), 32)
}

// ToPaddedBytes converts a big integer into a byte slice in big-endian order, ensuring the slice is at least n bytes long.
func ToPaddedBytes(bigint *big.Int, length int) []byte {
	if bigint.BitLen()/8 >= length {
		return bigint.Bytes()
	}
	result := make([]byte, length)
	FillBytes(bigint, result)
	return result
}

// FillBytes fills a byte slice with the absolute value of a big integer in big-endian order.
// The caller must ensure the buffer has enough space, otherwise the result will be truncated.
func FillBytes(bigint *big.Int, buffer []byte) {
	index := len(buffer)
	for _, digit := range bigint.Bits() {
		for j := 0; j < BytesPerWord && index > 0; j++ {
			index--
			buffer[index] = byte(digit)
			digit >>= 8
		}
	}
}

// ToSigned256 converts the input number to its 256-bit two's complement signed representation.
// The input must not exceed 256 bits.
func ToSigned256(x *big.Int) *big.Int {
	if x.Cmp(bigLimit255) < 0 {
		return x
	}
	return new(big.Int).Sub(x, bigLimit256)
}
