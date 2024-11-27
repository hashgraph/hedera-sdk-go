package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ContractLogInfo is the log info for events returned by a function
type ContractLogInfo struct {
	ContractID ContractID
	Bloom      []byte
	Topics     [][]byte
	Data       []byte
}

func _ContractLogInfoFromProtobuf(pb *services.ContractLoginfo) ContractLogInfo {
	if pb == nil {
		return ContractLogInfo{}
	}

	contractID := ContractID{}
	if pb.ContractID != nil {
		contractID = *_ContractIDFromProtobuf(pb.ContractID)
	}

	return ContractLogInfo{
		ContractID: contractID,
		Bloom:      pb.Bloom,
		Topics:     pb.Topic,
		Data:       pb.Data,
	}
}

func (logInfo ContractLogInfo) _ToProtobuf() *services.ContractLoginfo {
	return &services.ContractLoginfo{
		ContractID: logInfo.ContractID._ToProtobuf(),
		Bloom:      logInfo.Bloom,
		Topic:      logInfo.Topics,
		Data:       logInfo.Data,
	}
}
