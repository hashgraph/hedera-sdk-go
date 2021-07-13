package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

// ContractLogInfo is the log info for events returned by a function
type ContractLogInfo struct {
	ContractID ContractID
	Bloom      []byte
	Topics     [][]byte
	Data       []byte
}

func contractLogInfoFromProtobuf(pb *services.ContractLoginfo, networkName *NetworkName) ContractLogInfo {
	if pb == nil {
		return ContractLogInfo{}
	}
	return ContractLogInfo{
		ContractID: contractIDFromProtobuf(pb.ContractID, networkName),
		Bloom:      pb.Bloom,
		Topics:     pb.Topic,
		Data:       pb.Data,
	}
}

func (logInfo ContractLogInfo) toProtobuf() *services.ContractLoginfo {
	return &services.ContractLoginfo{
		ContractID: logInfo.ContractID.toProtobuf(),
		Bloom:      logInfo.Bloom,
		Topic:      logInfo.Topics,
		Data:       logInfo.Data,
	}
}
