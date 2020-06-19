package hedera

import (
	"encoding/binary"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ContractFunctionResult is the result returned by a call to a smart contract function. This is The response to
// a ContractCallQuery, and is in the record for a ContractCallQuery.
type ContractFunctionResult struct {
	// ContractID is the smart contract instance whose function was called
	ContractID *ContractID
	// ContractCallResult is the result returned by the function
	ContractCallResult []byte
	// ErrorMessage is the message returned in the case there was an error during smart contract execution
	ErrorMessage string
	// Bloom is the bloom filter for record
	Bloom []byte
	// GasUsed is the amount of gas used to execute the contract function
	GasUsed uint64
	// LogInfo is the log info for events returned by the function
	LogInfo []ContractLogInfo
}

// GetBool gets a solidity bool from the result at the given index
func (result ContractFunctionResult) GetBool(index uint64) bool {
	return result.GetUint32(index) == 1
}

// GetAddress gets a solidity address from the result at the given index
func (result ContractFunctionResult) GetAddress(index uint64) []byte {
	return result.ContractCallResult[(index*32)+12 : (index*32)+32]
}

// GetInt8 gets a solidity int8 from the result at the given index
func (result ContractFunctionResult) GetInt8(index uint64) int8 {
	return int8(result.ContractCallResult[index*32+31])
}

// GetInt32 gets a solidity int32 from the result at the given index
func (result ContractFunctionResult) GetInt32(index uint64) int32 {
	return int32(binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32]))
}

// GetInt64 gets a solidity int64 from the result at the given index
func (result ContractFunctionResult) GetInt64(index uint64) int64 {
	return int64(binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32]))
}

// GetInt256 gets a solidity int256 from the result at the given index
func (result ContractFunctionResult) GetInt256(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetUint8 gets a solidity uint8 from the result at the given index
func (result ContractFunctionResult) GetUint8(index uint64) uint8 {
	return result.ContractCallResult[index*32+31]
}

// GetUint32 gets a solidity uint32 from the result at the given index
func (result ContractFunctionResult) GetUint32(index uint64) uint32 {
	return binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32])
}

// GetUint64 gets a solidity uint64 from the result at the given index
func (result ContractFunctionResult) GetUint64(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

// GetUint256 gets a solidity uint256 from the result at the given index
func (result ContractFunctionResult) GetUint256(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetBytes32 gets a solidity bytes32 from the result at the given index
func (result ContractFunctionResult) GetBytes32(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetString gets a string from the result at the given index
func (result ContractFunctionResult) GetString(index uint64) string {
	return string(result.GetBytes(index))
}

// GetBytes gets a byte array from the result at the given index
func (result ContractFunctionResult) GetBytes(index uint64) []byte {
	offset := result.GetUint64(index)
	length := binary.BigEndian.Uint64(result.ContractCallResult[offset+24 : offset+32])
	return result.ContractCallResult[offset+32 : offset+32+length]
}

// AsBytes returns the raw bytes of the ContractCallResult
func (result ContractFunctionResult) AsBytes() []byte {
	return result.ContractCallResult
}

func contractFunctionResultFromProto(pb *proto.ContractFunctionResult) ContractFunctionResult {
	infos := make([]ContractLogInfo, len(pb.LogInfo))

	for i, info := range pb.LogInfo {
		infos[i] = contractLogInfoFromProto(info)
	}

	result := ContractFunctionResult{
		ContractID:         nil,
		ContractCallResult: pb.ContractCallResult,
		ErrorMessage:       pb.ErrorMessage,
		Bloom:              pb.Bloom,
		GasUsed:            pb.GasUsed,
		LogInfo:            infos,
	}

	if pb.ContractID != nil {
		contractID := contractIDFromProto(pb.ContractID)
		result.ContractID = &contractID
	}

	return result
}
