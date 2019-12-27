package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type ContractFunctionResult struct {
	ContractID         ContractID
	contractCallResult []byte
	ErrorMessage       string
	Bloom              []byte
	GasUsed            uint64
	Logs               []ContractLogInfo
}

func contractFunctionResultFromProto(pb *proto.ContractFunctionResult) ContractFunctionResult {
	logs := []ContractLogInfo{}

	for _, log := range pb.LogInfo {
		logs = append(logs, contractLogInfoFromProto(log))
	}

	return ContractFunctionResult{
		ContractID:         contractIDFromProto(pb.ContractID),
		contractCallResult: pb.ContractCallResult,
		ErrorMessage:       pb.ErrorMessage,
		Bloom:              pb.Bloom,
		GasUsed:            pb.GasUsed,
		Logs:               logs,
	}
}

func (result ContractFunctionResult) AsBytes() []byte {
	return result.contractCallResult
}
