package hiero

// SPDX-License-Identifier: Apache-2.0

var hardenedBit uint32 = 0x80000000

// Harden the index
func ToHardenedIndex(index uint32) uint32 {
	return index | hardenedBit
}

// Check if the index is hardened
func IsHardenedIndex(index uint32) bool {
	return (index & hardenedBit) != 0
}
