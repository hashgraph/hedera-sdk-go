package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// ContractLogInfo is the log info for events returned by a function
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
