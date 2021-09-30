package hedera

import (
	"golang.org/x/crypto/sha3"
)

type ContractFunctionSelector struct {
	function   *string
	params     string
	paramTypes []_Solidity
}

type _Solidity struct {
	ty    argument
	array bool
}

type argument string

const (
	aBool     argument = "bool"
	aString   argument = "string"
	aInt8     argument = "int8"
	aInt32    argument = "int32"
	aInt64    argument = "int64"
	aInt256   argument = "int256"
	aUint8    argument = "uint8"
	aUint32   argument = "uint32"
	aUint64   argument = "uint64"
	aUint256  argument = "uint256"
	aBytes    argument = "bytes"
	aBytes32  argument = "bytes32"
	aFunction argument = "function"
	aAddress  argument = "address"
)

func NewContractFunctionSelector(name string) ContractFunctionSelector {
	var function *string

	if name == "" {
		function = nil
	} else {
		function = &name
	}

	return ContractFunctionSelector{
		function:   function,
		params:     "",
		paramTypes: []_Solidity{},
	}
}

func (selector *ContractFunctionSelector) _AddParam(ty _Solidity) *ContractFunctionSelector {
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
	return selector._AddParam(_Solidity{
		ty:    aFunction,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddAddress() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aAddress,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddBool() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBool,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddString() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aString,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt8() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt8,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt32() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt32,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt64() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt64,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt256() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt256,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint8() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint8,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint32() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint32,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint64() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint64,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint256() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint256,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddBytes() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBytes,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddBytes32() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBytes32,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddAddressArray() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aAddress,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddBoolArray() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBool,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddStringArray() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aString,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddInt8Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt8,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddInt32Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt32,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddInt64Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt64,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddInt256Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt256,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddUint8Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint8,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddUint32Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint32,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddUint64Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint64,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddUint256Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint256,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddBytesArray() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBytes,
		array: true,
	})
}

func (selector *ContractFunctionSelector) AddBytes32Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBytes32,
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

func (selector *ContractFunctionSelector) _Build(function *string) []byte {
	if function != nil {
		selector.function = function
	} else if selector.function == nil {
		panic("unreachable: function name must be non-nil at this point")
	}

	hash := sha3.NewLegacyKeccak256()
	if _, err := hash.Write([]byte(selector.String())); err != nil {
		panic(err)
	}

	return hash.Sum(nil)[0:4]
}
