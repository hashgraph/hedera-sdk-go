package hedera

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
)

type ContractFunctionSelector struct {
	function   *string
	params     string
	paramTypes []solidity
}

type solidity struct {
	ty    argument
	array bool
}

type argument string

const (
	aBool     argument = "bool"
	aString   argument = "string"
	aInt32    argument = "int32"
	aInt64    argument = "int64"
	aInt256   argument = "int256"
	aBytes    argument = "bytes"
	aFunction argument = "function"
	aAddress  argument = "address"
)

func NewContractFunctionSelector(function *string) ContractFunctionSelector {
	return ContractFunctionSelector{
		function:   function,
		params:     "",
		paramTypes: []solidity{},
	}
}

func (selector *ContractFunctionSelector) addParam(ty solidity) *ContractFunctionSelector {
	if len(selector.paramTypes) > 0 {
		selector.params += ","
	}

	selector.params += string(ty.ty)
	if ty.array {
		selector.params += "[]"
	}

	selector.paramTypes = append(selector.paramTypes, ty)
	return selector
}

func (selector *ContractFunctionSelector) AddFunction() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aFunction,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddAddress() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aAddress,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddBool() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aBool,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddString() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aString,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt32() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aInt32,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt64() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aInt64,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt256() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aInt256,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddBytes() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aBytes,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddAddressArray() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aAddress,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddBoolArray() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aBool,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddStringArray() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aString,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddInt32Array() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aInt32,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddInt64Array() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aInt64,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddInt256Array() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aInt256,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddBytesArray() *ContractFunctionSelector {
	return selector.addParam(solidity{
		ty:    aBytes,
		array: true,
	})
}

func (selector *ContractFunctionSelector) String() string {
	function := ""
	if selector.function != nil {
		function = *selector.function
	}

	return function + "(" + selector.params + ")"
}

func (selector *ContractFunctionSelector) build(function *string) ([]byte, error) {
	if function != nil {
		selector.function = function
	} else if selector.function == nil {
		return nil, errors.Unwrap(fmt.Errorf("Address is required to be 40 characters"))
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(selector.String()))
	return hash.Sum(nil)[0:4], nil
}
