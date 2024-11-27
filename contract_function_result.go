package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ContractFunctionResult is a struct which allows users to convert between solidity and Go types, and is typically
// returned by `ContractCallQuery` and is present in the transaction records of `ContractExecuteTransaction`.
// Use the methods `Get<Type>()` to get a parameter. Not all solidity types
// are supported out of the box, but the most common types are. The larger variants
// of number types return just the bytes for the integer instead of converting to a big int type.
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
	ContractNonces       []*ContractNonceInfo
	SignerNonce          int64
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

// GetInt16 gets a _Solidity int16 from the result at the given index
func (result ContractFunctionResult) GetInt16(index uint64) int16 {
	return int16(binary.BigEndian.Uint16(result.ContractCallResult[index*32+30 : (index+1)*32]))
}

// GetInt24 gets a _Solidity int24 from the result at the given index
func (result ContractFunctionResult) GetInt24(index uint64) int32 {
	return int32(binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32]))
}

// GetInt40 gets a _Solidity int40 from the result at the given index
func (result ContractFunctionResult) GetInt40(index uint64) int64 {
	return int64(binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32]))
}

// GetInt48 gets a _Solidity int48 from the result at the given index
func (result ContractFunctionResult) GetInt48(index uint64) int64 {
	return int64(binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32]))
}

// GetInt56 gets a _Solidity int56 from the result at the given index
func (result ContractFunctionResult) GetInt56(index uint64) int64 {
	return int64(binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32]))
}

// GetInt32 gets a _Solidity int32 from the result at the given index
func (result ContractFunctionResult) GetInt32(index uint64) int32 {
	return int32(binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32]))
}

// GetInt64 gets a _Solidity int64 from the result at the given index
func (result ContractFunctionResult) GetInt64(index uint64) int64 {
	return int64(binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32]))
}

// GetInt72 gets a _Solidity int72 from the result at the given index
func (result ContractFunctionResult) GetInt72(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt80 gets a _Solidity int80 from the result at the given index
func (result ContractFunctionResult) GetInt80(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt88 gets a _Solidity int88 from the result at the given index
func (result ContractFunctionResult) GetInt88(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt96 gets a _Solidity int96 from the result at the given index
func (result ContractFunctionResult) GetInt96(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt104 gets a _Solidity int104 from the result at the given index
func (result ContractFunctionResult) GetInt104(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt112 gets a _Solidity int112 from the result at the given index
func (result ContractFunctionResult) GetInt112(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt120 gets a _Solidity int120 from the result at the given index
func (result ContractFunctionResult) GetInt120(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt128 gets a _Solidity int128 from the result at the given index
func (result ContractFunctionResult) GetInt128(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt136 gets a _Solidity int136 from the result at the given index
func (result ContractFunctionResult) GetInt136(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt144 gets a _Solidity int144 from the result at the given index
func (result ContractFunctionResult) GetInt144(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt152 gets a _Solidity int152 from the result at the given index
func (result ContractFunctionResult) GetInt152(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt160 gets a _Solidity int160 from the result at the given index
func (result ContractFunctionResult) GetInt160(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt168 gets a _Solidity int168 from the result at the given index
func (result ContractFunctionResult) GetInt168(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt176 gets a _Solidity int176 from the result at the given index
func (result ContractFunctionResult) GetInt176(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt184 gets a _Solidity int184 from the result at the given index
func (result ContractFunctionResult) GetInt184(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt192 gets a _Solidity int192 from the result at the given index
func (result ContractFunctionResult) GetInt192(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt200 gets a _Solidity int200 from the result at the given index
func (result ContractFunctionResult) GetInt200(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt208 gets a _Solidity int208 from the result at the given index
func (result ContractFunctionResult) GetInt208(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt216 gets a _Solidity int216 from the result at the given index
func (result ContractFunctionResult) GetInt216(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt224 gets a _Solidity int224 from the result at the given index
func (result ContractFunctionResult) GetInt224(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt232 gets a _Solidity int232 from the result at the given index
func (result ContractFunctionResult) GetInt232(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt240 gets a _Solidity int240 from the result at the given index
func (result ContractFunctionResult) GetInt240(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt248 gets a _Solidity int248 from the result at the given index
func (result ContractFunctionResult) GetInt248(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetInt256 gets a _Solidity int256 from the result at the given index
func (result ContractFunctionResult) GetInt256(index uint64) []byte {
	return result.ContractCallResult[index*32 : index*32+32]
}

// GetBigInt gets an _Solidity integer from the result at the given index and returns it as a big.Int
func (result ContractFunctionResult) GetBigInt(index uint64) *big.Int {
	value := new(big.Int).SetBytes(result.ContractCallResult[index*32 : index*32+32])
	fromTwosComplement := ToSigned256(value)
	return fromTwosComplement
}

// GetUint8 gets a _Solidity uint8 from the result at the given index
func (result ContractFunctionResult) GetUint8(index uint64) uint8 {
	return result.ContractCallResult[index*32+31]
}

// GetUint16 gets a _Solidity uint16 from the result at the given index
func (result ContractFunctionResult) GetUint16(index uint64) uint16 {
	return binary.BigEndian.Uint16(result.ContractCallResult[index*32+30 : (index+1)*32])
}

// GetUint24 gets a _Solidity uint24 from the result at the given index
func (result ContractFunctionResult) GetUint24(index uint64) uint32 {
	return binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32])
}

// GetUint32 gets a _Solidity uint32 from the result at the given index
func (result ContractFunctionResult) GetUint32(index uint64) uint32 {
	return binary.BigEndian.Uint32(result.ContractCallResult[index*32+28 : (index+1)*32])
}

// GetUint40 gets a _Solidity uint40 from the result at the given index
func (result ContractFunctionResult) GetUint40(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

// GetUint48 gets a _Solidity uint48 from the result at the given index
func (result ContractFunctionResult) GetUint48(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

// GetUint56 gets a _Solidity uint56 from the result at the given index
func (result ContractFunctionResult) GetUint56(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

// GetUint64 gets a _Solidity uint64 from the result at the given index
func (result ContractFunctionResult) GetUint64(index uint64) uint64 {
	return binary.BigEndian.Uint64(result.ContractCallResult[index*32+24 : (index+1)*32])
}

// GetUint72 gets a _Solidity uint72 from the result at the given index
func (result ContractFunctionResult) GetUint72(index uint64) []byte {
	return result.ContractCallResult[index*32+23 : (index+1)*32]
}

// GetUint80 gets a _Solidity uint80 from the result at the given index
func (result ContractFunctionResult) GetUint80(index uint64) []byte {
	return result.ContractCallResult[index*32+22 : (index+1)*32]
}

// GetUint88 gets a _Solidity uint88 from the result at the given index
func (result ContractFunctionResult) GetUint88(index uint64) []byte {
	return result.ContractCallResult[index*32+21 : (index+1)*32]
}

// GetUint96 gets a _Solidity uint96 from the result at the given index
func (result ContractFunctionResult) GetUint96(index uint64) []byte {
	return result.ContractCallResult[index*32+20 : (index+1)*32]
}

// GetUint104 gets a _Solidity uint104 from the result at the given index
func (result ContractFunctionResult) GetUint104(index uint64) []byte {
	return result.ContractCallResult[index*32+19 : (index+1)*32]
}

// GetUint112 gets a _Solidity uint112 from the result at the given index
func (result ContractFunctionResult) GetUint112(index uint64) []byte {
	return result.ContractCallResult[index*32+18 : (index+1)*32]
}

// GetUint120 gets a _Solidity uint120 from the result at the given index
func (result ContractFunctionResult) GetUint120(index uint64) []byte {
	return result.ContractCallResult[index*32+17 : (index+1)*32]
}

// GetUint128 gets a _Solidity uint128 from the result at the given index
func (result ContractFunctionResult) GetUint128(index uint64) []byte {
	return result.ContractCallResult[index*32+16 : (index+1)*32]
}

// GetUint136 gets a _Solidity uint136 from the result at the given index
func (result ContractFunctionResult) GetUint136(index uint64) []byte {
	return result.ContractCallResult[index*32+15 : (index+1)*32]
}

// GetUint144 gets a _Solidity uint144 from the result at the given index
func (result ContractFunctionResult) GetUint144(index uint64) []byte {
	return result.ContractCallResult[index*32+14 : (index+1)*32]
}

// GetUint152 gets a _Solidity uint152 from the result at the given index
func (result ContractFunctionResult) GetUint152(index uint64) []byte {
	return result.ContractCallResult[index*32+13 : (index+1)*32]
}

// GetUint160 gets a _Solidity uint160 from the result at the given index
func (result ContractFunctionResult) GetUint160(index uint64) []byte {
	return result.ContractCallResult[index*32+12 : (index+1)*32]
}

// GetUint168 gets a _Solidity uint168 from the result at the given index
func (result ContractFunctionResult) GetUint168(index uint64) []byte {
	return result.ContractCallResult[index*32+11 : (index+1)*32]
}

// GetUint176 gets a _Solidity uint176 from the result at the given index
func (result ContractFunctionResult) GetUint176(index uint64) []byte {
	return result.ContractCallResult[index*32+10 : (index+1)*32]
}

// GetUint184 gets a _Solidity uint184 from the result at the given index
func (result ContractFunctionResult) GetUint184(index uint64) []byte {
	return result.ContractCallResult[index*32+9 : (index+1)*32]
}

// GetUint192 gets a _Solidity uint192 from the result at the given index
func (result ContractFunctionResult) GetUint192(index uint64) []byte {
	return result.ContractCallResult[index*32+8 : (index+1)*32]
}

// GetUint200 gets a _Solidity uint200 from the result at the given index
func (result ContractFunctionResult) GetUint200(index uint64) []byte {
	return result.ContractCallResult[index*32+7 : (index+1)*32]
}

// GetUint208 gets a _Solidity uint208 from the result at the given index
func (result ContractFunctionResult) GetUint208(index uint64) []byte {
	return result.ContractCallResult[index*32+6 : (index+1)*32]
}

// GetUint216 gets a _Solidity uint216 from the result at the given index
func (result ContractFunctionResult) GetUint216(index uint64) []byte {
	return result.ContractCallResult[index*32+5 : (index+1)*32]
}

// GetUint224 gets a _Solidity uint224 from the result at the given index
func (result ContractFunctionResult) GetUint224(index uint64) []byte {
	return result.ContractCallResult[index*32+4 : (index+1)*32]
}

// GetUint232 gets a _Solidity uint232 from the result at the given index
func (result ContractFunctionResult) GetUint232(index uint64) []byte {
	return result.ContractCallResult[index*32+3 : (index+1)*32]
}

// GetUint240 gets a _Solidity uint240 from the result at the given index
func (result ContractFunctionResult) GetUint240(index uint64) []byte {
	return result.ContractCallResult[index*32+2 : (index+1)*32]
}

// GetUint248 gets a _Solidity uint248 from the result at the given index
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

// GetResult parses the result of a contract call based on the given types string and returns the result as an interface.
// The "types" string should specify the Ethereum Solidity type of the contract call output.
// This includes types like "uint256", "address", "bool", "string", "string[]", etc.
// The type provided must match the actual type of the data returned by the contract call,
// otherwise the function will fail to unpack and return an error.
// The method returns the parsed result encapsulated in an interface{},
// allowing flexibility to handle various types of contract call results.
// For correct usage, the caller should perform a type assertion on the returned interface{}
// to convert it into the appropriate go type.
func (result ContractFunctionResult) GetResult(types string) (interface{}, error) {
	def := fmt.Sprintf(`[{ "name" : "method", "type": "function", "outputs": [{ "type": "%s" }]}]`, types)
	abi := ABI{}
	err := abi.UnmarshalJSON([]byte(def))
	if err != nil {
		return nil, err
	}

	parsedResult, err := abi.Methods["method"].Decode(result.ContractCallResult)

	if err != nil {
		return nil, err
	}
	return parsedResult["0"], nil
}

func extractInt64OrZero(pb *services.ContractFunctionResult) int64 {
	if pb.GetSignerNonce() != nil {
		return pb.SignerNonce.Value
	}
	return 0
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

	var nonces []*ContractNonceInfo
	if len(pb.ContractNonces) > 0 {
		nonces = make([]*ContractNonceInfo, len(pb.ContractNonces))
		for i, nonce := range pb.ContractNonces {
			nonces[i] = _ContractNonceInfoFromProtobuf(nonce)
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
		ContractNonces:     nonces,
		SignerNonce:        extractInt64OrZero(pb),
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
		SignerNonce:        wrapperspb.Int64(result.SignerNonce),
	}
}

// ToBytes returns the protobuf encoded bytes of the ContractFunctionResult
func (result *ContractFunctionResult) ToBytes() []byte {
	data, err := protobuf.Marshal(result._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// ContractFunctionResultFromBytes returns a ContractFunctionResult from the protobuf encoded bytes of a ContractFunctionResult
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
