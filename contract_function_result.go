package hedera

import (
	"encoding/binary"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractFunctionResult struct {
	ContractID         ContractID
	ContractCallResult []byte
	ErrorMessage       string
	Bloom              []byte
	GasUsed            uint64
	LogInfo            []ContractLogInfo
}

type ContractLogInfo struct {
	ContractID ContractID
	Bloom      []byte
	Topic      [][]byte
	Data       []byte
}

func (result ContractFunctionResult) GetBool(index uint64) bool {
	return result.GetUint32(index) == 1
}

func (result ContractFunctionResult) GetAddress(index uint64) []byte {
	return result.GetBytes(index)[12:32]
}

func (result ContractFunctionResult) GetFunction(index uint64) ([]byte, []byte) {
	return result.GetBytes(index)[8:28], result.GetBytes(index)[28:32]
}

func (result ContractFunctionResult) GetInt256(index uint64) []byte {
	return result.GetBytes(index)
}

func (result ContractFunctionResult) GetString(index uint64) string {
	return string(result.GetBytes(index))
}

func (result ContractFunctionResult) GetBytes(index uint64) []byte {
	length := uint64(result.GetUint32(index))
	return result.ContractCallResult[(index+1)*32 : (index+1)*32+length]
}

func (result ContractFunctionResult) GetUint32(index uint64) uint32 {
	return binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32])
}

func (result ContractFunctionResult) GetUint64(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

func contractFunctionResultFromProto(pb *proto.ContractFunctionResult) ContractFunctionResult {
	infos := []ContractLogInfo{}
	for _, info := range pb.LogInfo {
		infos = append(infos, ContractLogInfo{
			ContractID: contractIDFromProto(info.ContractID),
			Bloom:      info.Bloom,
			Topic:      info.Topic,
			Data:       info.Data,
		})
	}

	return ContractFunctionResult{
		ContractID:         contractIDFromProto(pb.ContractID),
		ContractCallResult: pb.ContractCallResult,
		ErrorMessage:       pb.ErrorMessage,
		Bloom:              pb.Bloom,
		GasUsed:            pb.GasUsed,
		LogInfo:            infos,
	}
}
