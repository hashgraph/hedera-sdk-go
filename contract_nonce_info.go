package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ContractID is the ID for a Hiero smart contract
type ContractNonceInfo struct {
	ContractID *ContractID
	Nonce      int64
}

func _ContractNonceInfoFromProtobuf(contractNonceInfo *services.ContractNonceInfo) *ContractNonceInfo {
	if contractNonceInfo == nil {
		return nil
	}

	return &ContractNonceInfo{
		ContractID: _ContractIDFromProtobuf(contractNonceInfo.GetContractId()),
		Nonce:      contractNonceInfo.GetNonce(),
	}
}
