package hedera

import "github.com/hashgraph/hedera-sdk-go/v2/proto"

// ContractLogInfo is the log info for events returned by a function
type ContractLogInfo struct {
	ContractID ContractID
	Bloom      []byte
	Topics     [][]byte
	Data       []byte
}

func contractLogInfoFromProtobuf(pb *proto.ContractLoginfo) ContractLogInfo {
	if pb == nil {
		return ContractLogInfo{}
	}

	contractID := ContractID{}
	if pb.ContractID != nil {
		contractID = *contractIDFromProtobuf(pb.ContractID)
	}

	return ContractLogInfo{
		ContractID: contractID,
		Bloom:      pb.Bloom,
		Topics:     pb.Topic,
		Data:       pb.Data,
	}
}

func (logInfo ContractLogInfo) toProtobuf() *proto.ContractLoginfo {
	return &proto.ContractLoginfo{
		ContractID: logInfo.ContractID.toProtobuf(),
		Bloom:      logInfo.Bloom,
		Topic:      logInfo.Topics,
		Data:       logInfo.Data,
	}
}
