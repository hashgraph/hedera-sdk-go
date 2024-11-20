package hiero

// SPDX-License-Identifier: Apache-2.0

type ContractStateChange struct {
	ContractID     *ContractID
	StorageChanges []*StorageChange
}
