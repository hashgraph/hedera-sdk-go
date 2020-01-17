package hedera

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type ContractFunctionParams struct {
	function  ContractFunctionSelector
	arguments []Argument
}

type Argument struct {
	value   []byte
	dynamic bool
}

func NewContractFunctionParams() *ContractFunctionParams {
	return &ContractFunctionParams{
		function:  NewContractFunctionSelector(""),
		arguments: []Argument{},
	}
}

func (contract *ContractFunctionParams) AddBool(value bool) *ContractFunctionParams {
	argument := NewArgument()

	if value {
		argument.value[31] = 1
	} else {
		argument.value[31] = 0
	}

	contract.function.AddBool()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddFunction(address string, selector ContractFunctionSelector) (*ContractFunctionParams, error) {
	if len(address) != 40 {
		return contract, errors.Unwrap(fmt.Errorf("Address is required to be 40 characters"))
	}

	argument := NewArgument()
	argument.dynamic = false

	addressBytes, err := hex.DecodeString(address)
	if err != nil {
		return contract, err
	}

	bytes := make([]byte, 12)
	bytes = append(bytes, addressBytes[0:20]...)

	function, err := selector.build(nil)
	if err != nil {
		return contract, nil
	}

	bytes = append(bytes, function[0:4]...)
	argument.value = bytes

	contract.function.AddFunction()
	contract.arguments = append(contract.arguments, argument)
	return contract, nil
}

func (contract *ContractFunctionParams) AddInt8(value int8) *ContractFunctionParams {
	argument := NewArgument()

	argument.value[31] = uint8(value)

	contract.function.AddInt8()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddInt32(value int32) *ContractFunctionParams {
	argument := NewArgument()

	binary.BigEndian.PutUint32(argument.value[28:32], uint32(value))

	contract.function.AddInt32()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddInt64(value int64) *ContractFunctionParams {
	argument := NewArgument()

	binary.BigEndian.PutUint64(argument.value[24:32], uint64(value))

	contract.function.AddInt64()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddInt256(value []byte) *ContractFunctionParams {
	argument := NewArgument()

	argument.value = value

	contract.function.AddInt256()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddUint8(value uint8) *ContractFunctionParams {
	argument := NewArgument()

	argument.value[31] = value

	contract.function.AddUint8()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddUint32(value uint32) *ContractFunctionParams {
	argument := NewArgument()

	binary.BigEndian.PutUint32(argument.value[28:32], value)

	contract.function.AddUint32()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddUint64(value uint64) *ContractFunctionParams {
	argument := NewArgument()

	binary.BigEndian.PutUint64(argument.value[24:32], value)

	contract.function.AddUint64()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddUint256(value []byte) *ContractFunctionParams {
	argument := NewArgument()

	argument.value = value

	contract.function.AddUint256()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddInt32Array(value []int32) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	result := make([]byte, len(value)+32)

	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	for i, v := range value {
		binary.BigEndian.PutUint32(result[i*32+32+28:i*32+32+32], uint32(v))
	}

	argument.value = result

	contract.function.AddInt32Array()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddInt64Array(value []int64) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	result := make([]byte, len(value)+32)

	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	for i, v := range value {
		binary.BigEndian.PutUint64(result[i*32+32+24:i*32+32+32], uint64(v))
	}

	argument.value = result

	contract.function.AddInt64Array()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddInt256Array(value [][32]byte) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	result := make([]byte, len(value)+32)

	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	for i, v := range value {
		copy(result[i*32+32:i*32+32+32], v[0:32])
	}

	argument.value = result

	contract.function.AddInt256Array()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddUint32Array(value []uint32) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	result := make([]byte, len(value)+32)

	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	for i, v := range value {
		binary.BigEndian.PutUint32(result[i*32+32+28:i*32+32+32], v)
	}

	argument.value = result

	contract.function.AddUint32Array()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddUint64Array(value []uint64) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	result := make([]byte, len(value)+32)

	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	for i, v := range value {
		binary.BigEndian.PutUint64(result[i*32+32+24:i*32+32+32], v)
	}

	argument.value = result

	contract.function.AddUint64Array()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddUint256Array(value [][32]byte) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	result := make([]byte, len(value)+32)

	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	for i, v := range value {
		copy(result[i*32+32:i*32+32+32], v[0:32])
	}

	argument.value = result

	contract.function.AddUint256Array()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddAddressArray(value []string) (*ContractFunctionParams, error) {
	argument := NewArgument()
	argument.dynamic = true

	result := make([]byte, len(value)+32)

	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	for i, v := range value {
		if len(v) != 40 {
			return contract, errors.Unwrap(fmt.Errorf("Address is required to be 40 characters"))
		}

		addressBytes, err := hex.DecodeString(v)
		if err != nil {
			return contract, err
		}

		copy(result[i*32+32+12:i*32+32+32], addressBytes[0:20])
	}

	argument.value = result

	contract.function.AddAddressArray()
	contract.arguments = append(contract.arguments, argument)
	return contract, nil
}

func (contract *ContractFunctionParams) AddString(value string) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	bytes := []byte(value)
	binary.BigEndian.PutUint64(argument.value[24:32], uint64(len(bytes)))
	argument.value = append(argument.value, bytes...)
	argument.value = append(argument.value, make([]byte, 32-len(bytes)%32)...)

	contract.function.AddString()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddBytes(value []byte) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	binary.BigEndian.PutUint64(argument.value[24:32], uint64(len(value)))
	argument.value = append(argument.value, value...)
	argument.value = append(argument.value, make([]byte, uint64(32-len(value)%32))...)

	contract.function.AddBytes()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddBytes32(value []byte) *ContractFunctionParams {
	argument := NewArgument()

	argument.value = value[:]

	contract.function.AddBytes32()
	contract.arguments = append(contract.arguments, argument)

	return contract
}

func (contract *ContractFunctionParams) AddAddress(value string) (*ContractFunctionParams, error) {
	if len(value) != 40 {
		return contract, errors.Unwrap(fmt.Errorf("Address is required to be 40 characters"))
	}

	addressBytes, err := hex.DecodeString(value)
	if err != nil {
		return contract, err
	}

	argument := NewArgument()
	argument.dynamic = false

	bytes := make([]byte, 12)
	bytes = append(bytes, addressBytes...)

	argument.value = bytes

	contract.function.AddAddress()
	contract.arguments = append(contract.arguments, argument)
	return contract, nil
}

func (contract *ContractFunctionParams) AddBytesArray(value [][]byte) *ContractFunctionParams {
	argument := NewArgument()

	argument.dynamic = true
	argument.value = bytesArray(value)

	contract.function.AddBytesArray()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddByte32sArray(value [][]byte) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	result := make([]byte, len(value)+32)

	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	for i, v := range value {
		copy(result[i*32+32:i*32+32+32], v[0:32])
	}

	argument.value = result

	contract.function.AddBytes32Array()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) AddStringArray(value []string) *ContractFunctionParams {
	argument := NewArgument()
	argument.dynamic = true

	var bytes [][]byte
	for _, s := range value {
		bytes = append(bytes, []byte(s))
	}

	argument.value = bytesArray(bytes)
	contract.function.AddStringArray()
	contract.arguments = append(contract.arguments, argument)
	return contract
}

func (contract *ContractFunctionParams) build(functionName *string) ([]byte, error) {
	function, err := contract.function.build(functionName)
	if err != nil {
		return []byte{}, nil
	}

	length := uint64(0)

	functionOffset := uint64(0)
	if functionName != nil {
		functionOffset = uint64(4)
	}

	for _, argument := range contract.arguments {
		length += uint64(32)
		if argument.dynamic {
			length += uint64(len(argument.value))
		}
	}

	result := make([]byte, length+functionOffset)
	copy(result[0:4], function)

	offset := uint64(len(contract.arguments) * 32)

	for i, argument := range contract.arguments {
		j := uint64(i)
		if argument.dynamic {
			binary.BigEndian.PutUint64(result[(j*32+functionOffset)+24:(j+1)*32+functionOffset], offset)
			copy(result[offset+functionOffset:], argument.value)
			offset += uint64(len(argument.value))
		} else {
			copy(result[j*32+functionOffset:((j+1)*32)+functionOffset], argument.value)
		}
	}

	return result, nil
}

func NewArgument() Argument {
	return Argument{
		value:   make([]byte, 32),
		dynamic: false,
	}
}

func bytesArray(value [][]byte) []byte {
	// Calculate Length of final result
	length := uint64(0)
	for _, s := range value {
		length += 32 + 32
		sbytes := []byte(s)
		if len(sbytes)/32 == 0 {
			length += 32
		} else {
			length += uint64(((len(sbytes) / 32) + 1) * 32)
		}
	}

	// Zero initialize final resulting byte array
	result := make([]byte, length+32)

	// Write length of array into the first 32 bytes
	binary.BigEndian.PutUint64(result[24:32], uint64(len(value)))

	// Create array of byte arrays to hold each string value
	// Needed to concat later
	arguments := make([][]byte, len(value))

	// Convert each argument into bytes, and push each argument
	// into the argument list
	for i, s := range value {
		// Get the length of the current argument (again)
		length = 0
		if len(s)/32 == 0 {
			length = 32
		} else {
			length = uint64(((len(s) / 32) + 1) * 32)
		}

		// Create byte array of correct size
		// Length of value to the nearest 32 byte boundry +
		// 32 bytes to store the length
		bytes := make([]byte, length+32)

		// Write length into first 32 bytes
		binary.BigEndian.PutUint64(bytes[24:32], uint64(len(s)))

		// Copy string as bytes to the rest of the buffer
		copy(bytes[32:], s)

		// Set the argument bytes to be used later
		arguments[i] = bytes
	}

	// Initialize offset to the number of strings
	offset := uint64(len(value) * 32)

	// For each argument, write the offset into result
	// and the argument value (which includes data and length already)
	for i, s := range arguments {
		binary.BigEndian.PutUint64(result[(i+1)*32+24:(i+2)*32], offset)
		copy(result[offset+32:offset+32+uint64(len(s))], s)
		offset += uint64(len(s))
	}

	return result
}
