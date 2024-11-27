//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeBytesBound(t *testing.T) {
	typ, _ := NewType("tuple(string)")
	decodeTuple(typ, nil) // it should not panic
}

func TestDecodeDynamicLengthOutOfBounds(t *testing.T) {
	input := []byte("00000000000000000000000000000000\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00 00000000000000000000000000")
	typ, _ := NewType("tuple(bytes32, bytes, bytes)")

	_, err := Decode(typ, input)
	require.Error(t, err)
}
