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

func (result ContractFunctionResult) GetBool(index uint64) bool {
	return result.GetUint32(index) == 1
}

func (result ContractFunctionResult) GetAddress(index uint64) []byte {
	return result.ContractCallResult[(index*32)+12 : (index*32)+32]
}

func (result ContractFunctionResult) GetInt32(index uint64) int32 {
	return int32(binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32]))
}

func (result ContractFunctionResult) GetInt64(index uint64) int64 {
	return int64(binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32]))
}

func (result ContractFunctionResult) GetInt256(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

func (result ContractFunctionResult) GetUint32(index uint64) uint32 {
	return binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32])
}

func (result ContractFunctionResult) GetUint64(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

func (result ContractFunctionResult) GetUint256(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

func (result ContractFunctionResult) GetString(index uint64) string {
	return string(result.GetBytes(index))
}

func (result ContractFunctionResult) GetBytes(index uint64) []byte {
	offset := result.GetUint64(index)
	length := binary.BigEndian.Uint64(result.ContractCallResult[offset+24 : offset+32])
	return result.ContractCallResult[offset+32 : offset+32+length]
}

func (result ContractFunctionResult) AsBytes() []byte {
	return result.ContractCallResult
}

func contractFunctionResultFromProto(pb *proto.ContractFunctionResult) ContractFunctionResult {
	infos := []ContractLogInfo{}
	for _, info := range pb.LogInfo {
		infos = append(infos, ContractLogInfo{
			ContractID: contractIDFromProto(info.ContractID),
			Bloom:      info.Bloom,
			Topics:     info.Topic,
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
