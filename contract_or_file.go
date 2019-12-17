package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type ContractIdOrFileId interface {
	toProtoContractIdOrFile() (*proto.FileID, *proto.ContractID, int32)
}

func (id FileID) toProtoContractIdOrFile() (*proto.FileID, *proto.ContractID, int32) {
	return id.toProto(), &proto.ContractID{}, 0
}

func (id ContractID) toProtoContractIdOrFile() (*proto.FileID, *proto.ContractID, int32) {
	return &proto.FileID{}, id.toProto(), 1
}
