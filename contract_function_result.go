package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"encoding/binary"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ContractFunctionResult is a struct which allows users to convert between solidity and Go types, and is typically
// returned by `ContractCallQuery` and is present in the transaction records of `ContractExecuteTransaction`.
// Use the methods `Get<Type>()` to get a parameter. Not all solidity types
// are supported out of the box, but the most common types are. The larger variants
// of number types return just the bytes for the integer instead of converting to a big int type.
// To convert those bytes into a usable integer using "github.com/ethereum/go-ethereum/common/math" and "math/big" do the following:
// ```
// contractFunctionResult.GetUint256(<index>)
// bInt := new(big.Int)
// bInt.SetBytes(query.GetUint256(0))
// ```
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
	// Deprecated
	CreatedContractIDs []ContractID
	// Deprecated
	ContractStateChanges []ContractStateChange
	EvmAddress           ContractID
	GasAvailable         int64
	Amount               Hbar
	FunctionParameters   []byte
}

// GetBool gets a _Solidity bool from the result at the given index
func (result ContractFunctionResult) GetBool(index uint64) bool {
	return result.GetUint32(index) == 1
}

// GetAddress gets a _Solidity address from the result at the given index
func (result ContractFunctionResult) GetAddress(index uint64) []byte {
	return result.ContractCallResult[(index*32)+12 : (index*32)+32]
}

// GetInt8 gets a _Solidity int8 from the result at the given index
func (result ContractFunctionResult) GetInt8(index uint64) int8 {
	return int8(result.ContractCallResult[index*32+31])
}

func (result ContractFunctionResult) GetInt16(index uint64) uint16 {
	return binary.BigEndian.Uint16(result.ContractCallResult[index*32+30 : (index+1)*32])
}

func (result ContractFunctionResult) GetInt24(index uint64) uint32 {
	return binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32])
}

func (result ContractFunctionResult) GetInt40(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

func (result ContractFunctionResult) GetInt48(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

func (result ContractFunctionResult) GetInt56(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

// GetInt32 gets a _Solidity int32 from the result at the given index
func (result ContractFunctionResult) GetInt32(index uint64) int32 {
	return int32(binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32]))
}

// GetInt64 gets a _Solidity int64 from the result at the given index
func (result ContractFunctionResult) GetInt64(index uint64) int64 {
	return int64(binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32]))
}

func (result ContractFunctionResult) GetInt72(index uint64) []byte {
	return result.ContractCallResult[index*32+23 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt80(index uint64) []byte {
	return result.ContractCallResult[index*32+22 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt88(index uint64) []byte {
	return result.ContractCallResult[index*32+21 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt96(index uint64) []byte {
	return result.ContractCallResult[index*32+20 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt104(index uint64) []byte {
	return result.ContractCallResult[index*32+19 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt112(index uint64) []byte {
	return result.ContractCallResult[index*32+18 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt120(index uint64) []byte {
	return result.ContractCallResult[index*32+17 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt128(index uint64) []byte {
	return result.ContractCallResult[index*32+16 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt136(index uint64) []byte {
	return result.ContractCallResult[index*32+15 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt144(index uint64) []byte {
	return result.ContractCallResult[index*32+14 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt152(index uint64) []byte {
	return result.ContractCallResult[index*32+13 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt160(index uint64) []byte {
	return result.ContractCallResult[index*32+12 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt168(index uint64) []byte {
	return result.ContractCallResult[index*32+11 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt176(index uint64) []byte {
	return result.ContractCallResult[index*32+10 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt184(index uint64) []byte {
	return result.ContractCallResult[index*32+9 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt192(index uint64) []byte {
	return result.ContractCallResult[index*32+8 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt200(index uint64) []byte {
	return result.ContractCallResult[index*32+7 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt208(index uint64) []byte {
	return result.ContractCallResult[index*32+6 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt216(index uint64) []byte {
	return result.ContractCallResult[index*32+5 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt224(index uint64) []byte {
	return result.ContractCallResult[index*32+4 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt232(index uint64) []byte {
	return result.ContractCallResult[index*32+3 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt240(index uint64) []byte {
	return result.ContractCallResult[index*32+2 : (index+1)*32]
}

func (result ContractFunctionResult) GetInt248(index uint64) []byte {
	return result.ContractCallResult[index*32+1 : (index+1)*32]
}

// GetInt256 gets a _Solidity int256 from the result at the given index
func (result ContractFunctionResult) GetInt256(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetUint8 gets a _Solidity uint8 from the result at the given index
func (result ContractFunctionResult) GetUint8(index uint64) uint8 {
	return result.ContractCallResult[index*32+31]
}

func (result ContractFunctionResult) GetUint16(index uint64) uint16 {
	return binary.BigEndian.Uint16(result.ContractCallResult[index*32+30 : (index+1)*32])
}

func (result ContractFunctionResult) GetUint24(index uint64) uint32 {
	return binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32])
}

// GetUint32 gets a _Solidity uint32 from the result at the given index
func (result ContractFunctionResult) GetUint32(index uint64) uint32 {
	return binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32])
}

func (result ContractFunctionResult) GetUint40(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

func (result ContractFunctionResult) GetUint48(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

func (result ContractFunctionResult) GetUint56(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

// GetUint64 gets a _Solidity uint64 from the result at the given index
func (result ContractFunctionResult) GetUint64(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

func (result ContractFunctionResult) GetUint72(index uint64) []byte {
	return result.ContractCallResult[index*32+23 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint80(index uint64) []byte {
	return result.ContractCallResult[index*32+22 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint88(index uint64) []byte {
	return result.ContractCallResult[index*32+21 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint96(index uint64) []byte {
	return result.ContractCallResult[index*32+20 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint104(index uint64) []byte {
	return result.ContractCallResult[index*32+19 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint112(index uint64) []byte {
	return result.ContractCallResult[index*32+18 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint120(index uint64) []byte {
	return result.ContractCallResult[index*32+17 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint128(index uint64) []byte {
	return result.ContractCallResult[index*32+16 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint136(index uint64) []byte {
	return result.ContractCallResult[index*32+15 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint144(index uint64) []byte {
	return result.ContractCallResult[index*32+14 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint152(index uint64) []byte {
	return result.ContractCallResult[index*32+13 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint160(index uint64) []byte {
	return result.ContractCallResult[index*32+12 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint168(index uint64) []byte {
	return result.ContractCallResult[index*32+11 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint176(index uint64) []byte {
	return result.ContractCallResult[index*32+10 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint184(index uint64) []byte {
	return result.ContractCallResult[index*32+9 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint192(index uint64) []byte {
	return result.ContractCallResult[index*32+8 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint200(index uint64) []byte {
	return result.ContractCallResult[index*32+7 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint208(index uint64) []byte {
	return result.ContractCallResult[index*32+6 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint216(index uint64) []byte {
	return result.ContractCallResult[index*32+5 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint224(index uint64) []byte {
	return result.ContractCallResult[index*32+4 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint232(index uint64) []byte {
	return result.ContractCallResult[index*32+3 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint240(index uint64) []byte {
	return result.ContractCallResult[index*32+2 : (index+1)*32]
}

func (result ContractFunctionResult) GetUint248(index uint64) []byte {
	return result.ContractCallResult[index*32+1 : (index+1)*32]
}

// GetUint256 gets a _Solidity uint256 from the result at the given index
func (result ContractFunctionResult) GetUint256(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetBytes32 gets a _Solidity bytes32 from the result at the given index
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

func _ContractFunctionResultFromProtobuf(pb *services.ContractFunctionResult) ContractFunctionResult {
	infos := make([]ContractLogInfo, len(pb.LogInfo))

	for i, info := range pb.LogInfo {
		infos[i] = _ContractLogInfoFromProtobuf(info)
	}

	createdContractIDs := make([]ContractID, 0)
	for _, id := range pb.CreatedContractIDs { // nolint
		temp := _ContractIDFromProtobuf(id)
		if temp != nil {
			createdContractIDs = append(createdContractIDs, *temp)
		}
	}

	var evm ContractID
	if len(pb.EvmAddress.GetValue()) > 0 {
		evm = ContractID{
			Shard:      0,
			Realm:      0,
			Contract:   0,
			EvmAddress: pb.EvmAddress.GetValue(),
			checksum:   nil,
		}
	}

	result := ContractFunctionResult{
		ContractCallResult: pb.ContractCallResult,
		ErrorMessage:       pb.ErrorMessage,
		Bloom:              pb.Bloom,
		GasUsed:            pb.GasUsed,
		LogInfo:            infos,
		CreatedContractIDs: createdContractIDs,
		EvmAddress:         evm,
		GasAvailable:       pb.Gas,
		Amount:             HbarFromTinybar(pb.Amount),
		FunctionParameters: pb.FunctionParameters,
	}

	if pb.ContractID != nil {
		result.ContractID = _ContractIDFromProtobuf(pb.ContractID)
	}

	return result
}

func (result ContractFunctionResult) _ToProtobuf() *services.ContractFunctionResult {
	infos := make([]*services.ContractLoginfo, len(result.LogInfo))

	for i, info := range result.LogInfo {
		infos[i] = info._ToProtobuf()
	}

	return &services.ContractFunctionResult{
		ContractID:         result.ContractID._ToProtobuf(),
		ContractCallResult: result.ContractCallResult,
		ErrorMessage:       result.ErrorMessage,
		Bloom:              result.Bloom,
		GasUsed:            result.GasUsed,
		LogInfo:            infos,
		EvmAddress:         &wrapperspb.BytesValue{Value: result.EvmAddress.EvmAddress},
		Gas:                result.GasAvailable,
		Amount:             result.Amount.AsTinybar(),
		FunctionParameters: result.FunctionParameters,
	}
}

func (result *ContractFunctionResult) ToBytes() []byte {
	data, err := protobuf.Marshal(result._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func ContractFunctionResultFromBytes(data []byte) (ContractFunctionResult, error) {
	if data == nil {
		return ContractFunctionResult{}, errByteArrayNull
	}
	pb := services.ContractFunctionResult{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ContractFunctionResult{}, err
	}

	return _ContractFunctionResultFromProtobuf(&pb), nil
}
