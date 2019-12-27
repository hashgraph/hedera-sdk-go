package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type ContractLogInfo struct {
	ContractID ContractID
	Bloom      []byte
	Topics     [][]byte
	data       []byte
}

func contractLogInfoFromProto(pb *proto.ContractLoginfo) ContractLogInfo {
	return ContractLogInfo{
		ContractID: contractIDFromProto(pb.ContractID),
		Bloom:      pb.Bloom,
		Topics:     pb.Topic,
		data:       pb.Data,
	}
}
