package hiero

import (
	"math/big"
)

// SPDX-License-Identifier: Apache-2.0

type StorageChange struct {
	Slot         *big.Int
	ValueRead    *big.Int
	ValueWritten *big.Int
}

// ToBytes returns the byte representation of the StorageChange
func (storageChange *StorageChange) ToBytes() []byte {
	return []byte{}
}

// StorageChangeFromBytes returns a StorageChange from a byte array
func StorageChangeFromBytes(data []byte) (StorageChange, error) {
	return StorageChange{}, nil
}
