package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type ContractIdOrFileID interface {
	toProtoContractIDOrFile() (*proto.FileID, *proto.ContractID, int32)
}

func (id FileID) toProtoContractIDOrFile() (*proto.FileID, *proto.ContractID, int32) {
	return id.toProto(), &proto.ContractID{}, 0
}

func (id ContractID) toProtoContractIDOrFile() (*proto.FileID, *proto.ContractID, int32) {
	return &proto.FileID{}, id.toProto(), 1
}
