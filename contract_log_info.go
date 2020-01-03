package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type ContractLogInfo struct {
	ContractID ContractID
	Bloom      []byte
	Topics     [][]byte
	Data       []byte
}

func contractLogInfoFromProto(pb *proto.ContractLoginfo) ContractLogInfo {
	return ContractLogInfo{
		ContractID: contractIDFromProto(pb.ContractID),
		Bloom:      pb.Bloom,
		Topics:     pb.Topic,
		Data:       pb.Data,
	}
}
