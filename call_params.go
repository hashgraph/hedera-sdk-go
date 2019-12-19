package hedera

import (
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
)

type CallParams struct {
	function  FunctionSelector
	arguments []Argument
}

type FunctionSelector struct {
	function   string
	needsComma bool
}

type Argument struct {
	value   []byte
	dynamic bool
}

func NewCallParams(function *string) *CallParams {
	return &CallParams{
		function:  NewFunctionSelector(function),
		arguments: []Argument{},
	}
}

func (call *CallParams) SetFunction(function string) *CallParams {
	call.function.SetFunction(function)
	return call
}

func (call *CallParams) AddBool(value bool) *CallParams {
	argument := NewArgument()

	if value {
		argument.value[31] = 1
	} else {
		argument.value[31] = 0
	}

	call.function.AddParamType("bool")
	call.arguments = append(call.arguments, argument)

	return call
}

func (call *CallParams) AddInt(value interface{}) (*CallParams, error) {
	argument := NewArgument()

	switch Value := value.(type) {
	case uint8:
		argument.value[31] = Value
	case int8:
		argument.value[31] = uint8(Value)
	case uint16:
		binary.BigEndian.PutUint16(argument.value[30:32], Value)
	case int16:
		binary.BigEndian.PutUint16(argument.value[30:32], uint16(Value))
	case uint32:
		binary.BigEndian.PutUint32(argument.value[28:32], Value)
	case int32:
		binary.BigEndian.PutUint32(argument.value[28:32], uint32(Value))
	case uint64:
		binary.BigEndian.PutUint64(argument.value[24:32], Value)
	case int64:
		binary.BigEndian.PutUint64(argument.value[24:32], uint64(Value))
	default:
		return call, errors.Unwrap(fmt.Errorf("Expected int variant, received %v", reflect.TypeOf(value).String()))
	}

	call.function.AddParamType(reflect.TypeOf(value).String())
	call.arguments = append(call.arguments, argument)

	return call, nil
}

func (call *CallParams) AddString(value string) *CallParams {
	argument := NewArgument()
	argument.dynamic = true

	bytes := []byte(value)
	binary.BigEndian.PutUint64(argument.value[24:32], uint64(len(bytes)))
	argument.value = append(argument.value, bytes...)
	argument.value = append(argument.value, make([]byte, 32-len(bytes)%32)...)

	call.function.AddParamType("string")
	call.arguments = append(call.arguments, argument)
	return call
}

func (call *CallParams) AddBytes(value []byte) *CallParams {
	argument := NewArgument()
	argument.dynamic = true

	binary.BigEndian.PutUint64(argument.value[24:32], uint64(len(value)))
	argument.value = append(argument.value, value...)
	argument.value = append(argument.value, make([]byte, uint64(32-len(value)%32))...)

	call.function.AddParamType("bytes")
	call.arguments = append(call.arguments, argument)
	return call
}

func (call *CallParams) AddBytesFixed(value []byte, length uint) (*CallParams, error) {
	argument := NewArgument()
	argument.dynamic = false

	if length > 32 || len(value) > 32 {
		return call, errors.Unwrap(fmt.Errorf("Fixed bytes length cannot exceed 32 bytes"))
	}

	call.function.AddParamType(fmt.Sprintf("bytes%v", length))

	value = append(value, make([]byte, 32-len(value)%32)...)
	argument.value = value

	call.arguments = append(call.arguments, argument)
	return call, nil
}

func (call *CallParams) AddAddress(value []byte) (*CallParams, error) {
	if len(value) != 20 {
		return call, errors.Unwrap(fmt.Errorf("Address is required to be 20 bytes"))
	}

	argument := NewArgument()
	argument.dynamic = false

	bytes := make([]byte, 12)
	bytes = append(bytes, value...)

	argument.value = bytes

	call.function.AddParamType("address")
	call.arguments = append(call.arguments, argument)
	return call, nil
}

func (call *CallParams) AddFunction(address []byte, selector FunctionSelector) (*CallParams, error) {
	if len(address) != 20 {
		return call, errors.Unwrap(fmt.Errorf("Address is required to be 20 bytes"))
	}

	argument := NewArgument()
	argument.dynamic = false

	address = append(address, selector.Finish()...)

	bytes := make([]byte, 8)
	bytes = append(bytes, address...)

	argument.value = bytes

	call.function.AddParamType("function")
	call.arguments = append(call.arguments, argument)
	return call, nil
}

func (call *CallParams) AddStringArray(value []string) *CallParams {
	argument := NewArgument()
	argument.dynamic = true

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
		// Convert string to byte array
		sbytes := []byte(s)

		// Get the length of the current argument (again)
		length = 0
		if len(s)/32 == 0 {
			length = 32
		} else {
			length = uint64(((len(sbytes) / 32) + 1) * 32)
		}

		// Create byte array of correct size
		// Length of value to the nearest 32 byte boundry +
		// 32 bytes to store the length
		bytes := make([]byte, length+32)

		// Write length into first 32 bytes
		binary.BigEndian.PutUint64(bytes[24:32], uint64(len(sbytes)))

		// Copy string as bytes to the rest of the buffer
		copy(bytes[32:], sbytes)

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

	// Result should now contain the offests + each arguement
	argument.value = result

	call.function.AddParamType("string[]")
	call.arguments = append(call.arguments, argument)
	return call
}

func (call *CallParams) Finish() []byte {
	function := call.function.Finish()
	length := 0

	for _, argument := range call.arguments {
		length += 32
		if argument.dynamic {
			length += len(argument.value)
		}
	}

	result := make([]byte, length+4)
	copy(result[0:4], function)

	offset := uint64(len(call.arguments) * 32)

	for i, argument := range call.arguments {
		if argument.dynamic {
			binary.BigEndian.PutUint64(result[((i*32)+4)+24:(i+1)*32+4], offset)
			copy(result[offset+4:], argument.value)
			offset += uint64(len(argument.value))
		} else {
			copy(result[(i*32)+4:((i+1)*32)+4], argument.value)
		}
	}

	return result
}

func NewFunctionSelector(function *string) FunctionSelector {
	if function == nil {
		return FunctionSelector{
			function:   "",
			needsComma: false,
		}
	} else {
		return FunctionSelector{
			function:   *function,
			needsComma: false,
		}
	}
}

func (selector *FunctionSelector) SetFunction(function string) *FunctionSelector {
	selector.function = function + "("
	return selector
}

func (selector *FunctionSelector) AddParamType(ty string) *FunctionSelector {
	if selector.needsComma {
		selector.function += ","
	}
	selector.needsComma = true
	selector.function += ty
	return selector
}

func (selector *FunctionSelector) Finish() []byte {
	selector.function += ")"
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(selector.function))
	return hash.Sum(nil)[0:4]
}

func NewArgument() Argument {
	return Argument{
		value:   make([]byte, 32),
		dynamic: false,
	}
}
