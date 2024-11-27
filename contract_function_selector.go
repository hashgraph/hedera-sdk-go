package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"golang.org/x/crypto/sha3"
)

// A selector for a function with a given name.
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
	aInt16    argument = "int16"
	aInt24    argument = "int24"
	aInt32    argument = "int32"
	aInt40    argument = "int40"
	aInt48    argument = "int48"
	aInt56    argument = "int56"
	aInt64    argument = "int64"
	aInt72    argument = "int72"
	aInt80    argument = "int80"
	aInt88    argument = "int88"
	aInt96    argument = "int96"
	aInt104   argument = "int104"
	aInt112   argument = "int112"
	aInt120   argument = "int120"
	aInt128   argument = "int128"
	aInt136   argument = "int136"
	aInt144   argument = "int144"
	aInt152   argument = "int152"
	aInt160   argument = "int160"
	aInt168   argument = "int168"
	aInt176   argument = "int176"
	aInt184   argument = "int184"
	aInt192   argument = "int192"
	aInt200   argument = "int200"
	aInt208   argument = "int208"
	aInt216   argument = "int216"
	aInt224   argument = "int224"
	aInt232   argument = "int232"
	aInt240   argument = "int240"
	aInt248   argument = "int248"
	aInt256   argument = "int256"
	aUint8    argument = "uint8"
	aUint16   argument = "uint16"
	aUint24   argument = "uint24"
	aUint32   argument = "uint32"
	aUint40   argument = "uint40"
	aUint48   argument = "uint48"
	aUint56   argument = "uint56"
	aUint64   argument = "uint64"
	aUint72   argument = "uint72"
	aUint80   argument = "uint80"
	aUint88   argument = "uint88"
	aUint96   argument = "uint96"
	aUint104  argument = "uint104"
	aUint112  argument = "uint112"
	aUint120  argument = "uint120"
	aUint128  argument = "uint128"
	aUint136  argument = "uint136"
	aUint144  argument = "uint144"
	aUint152  argument = "uint152"
	aUint160  argument = "uint160"
	aUint168  argument = "uint168"
	aUint176  argument = "uint176"
	aUint184  argument = "uint184"
	aUint192  argument = "uint192"
	aUint200  argument = "uint200"
	aUint208  argument = "uint208"
	aUint216  argument = "uint216"
	aUint224  argument = "uint224"
	aUint232  argument = "uint232"
	aUint240  argument = "uint240"
	aUint248  argument = "uint248"
	aUint256  argument = "uint256"
	aBytes    argument = "bytes"
	aBytes32  argument = "bytes32"
	aFunction argument = "function"
	aAddress  argument = "address"
)

// NewContractFunctionSelector starts building a selector for a function with a given name.
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

// AddParam adds a parameter to the selector.
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

// AddFunction adds a function parameter to the selector.
func (selector *ContractFunctionSelector) AddFunction() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aFunction,
		array: false,
	})
}

// AddAddress adds an address parameter to the selector.
func (selector *ContractFunctionSelector) AddAddress() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aAddress,
		array: false,
	})
}

// AddBool adds a bool parameter to the selector.
func (selector *ContractFunctionSelector) AddBool() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBool,
		array: false,
	})
}

// AddString adds a string parameter to the selector.
func (selector *ContractFunctionSelector) AddString() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aString,
		array: false,
	})
}

// AddInt8 adds an int8 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt8() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt8,
		array: false,
	})
}

// AddInt16 adds an int16 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt16() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt16,
		array: false,
	})
}

// AddInt24 adds an int24 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt24() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt24,
		array: false,
	})
}

// AddInt32 adds an int32 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt32() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt32,
		array: false,
	})
}

// AddInt40 adds an int40 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt40() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt40,
		array: false,
	})
}

// AddInt48 adds an int48 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt48() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt48,
		array: false,
	})
}

// AddInt56 adds an int56 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt56() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt56,
		array: false,
	})
}

// AddInt64 adds an int64 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt64() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt64,
		array: false,
	})
}

// AddInt72 adds an int72 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt72() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt72,
		array: false,
	})
}

// AddInt80 adds an int80 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt80() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt80,
		array: false,
	})
}

// AddInt88 adds an int88 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt88() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt88,
		array: false,
	})
}

// AddInt96 adds an int96 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt96() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt96,
		array: false,
	})
}

// AddInt104 adds an int104 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt104() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt104,
		array: false,
	})
}

// AddInt112 adds an int112 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt112() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt112,
		array: false,
	})
}

// AddInt120 adds an int120 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt120() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt120,
		array: false,
	})
}

// AddInt128 adds an int128 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt128() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt128,
		array: false,
	})
}

// AddInt136 adds an int136 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt136() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt136,
		array: false,
	})
}

// AddInt144 adds an int144 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt144() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt144,
		array: false,
	})
}

// AddInt152 adds an int152 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt152() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt152,
		array: false,
	})
}

// AddInt160 adds an int160 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt160() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt160,
		array: false,
	})
}

// AddInt168 adds an int168 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt168() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt168,
		array: false,
	})
}

// AddInt176 adds an int176 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt176() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt176,
		array: false,
	})
}

// AddInt184 adds an int184 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt184() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt184,
		array: false,
	})
}

// AddInt192 adds an int192 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt192() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt192,
		array: false,
	})
}

// AddInt200 adds an int200 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt200() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt200,
		array: false,
	})
}

// AddInt208 adds an int208 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt208() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt208,
		array: false,
	})
}

// AddInt216 adds an int216 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt216() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt216,
		array: false,
	})
}

// AddInt224 adds an int224 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt224() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt224,
		array: false,
	})
}

// AddInt232 adds an int232 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt232() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt232,
		array: false,
	})
}

// AddInt240 adds an int240 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt240() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt240,
		array: false,
	})
}

// AddInt248 adds an int248 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt248() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt248,
		array: false,
	})
}

// AddInt256 adds an int256 parameter to the selector.
func (selector *ContractFunctionSelector) AddInt256() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt256,
		array: false,
	})
}

// AddUint8 adds a uint8 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint8() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint8,
		array: false,
	})
}

// AddUint16 adds a uint16 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint16() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint16,
		array: false,
	})
}

// AddUint24 adds a uint24 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint24() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint24,
		array: false,
	})
}

// AddUint32 adds a uint32 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint32() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint32,
		array: false,
	})
}

// AddUint40 adds a uint40 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint40() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint40,
		array: false,
	})
}

// AddUint48 adds a uint48 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint48() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint48,
		array: false,
	})
}

// AddUint56 adds a uint56 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint56() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint56,
		array: false,
	})
}

// AddUint64 adds a uint64 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint64() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint64,
		array: false,
	})
}

// AddUint72 adds a uint72 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint72() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint72,
		array: false,
	})
}

// AddUint80 adds a uint80 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint80() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint80,
		array: false,
	})
}

// AddUint88 adds a uint88 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint88() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint88,
		array: false,
	})
}

// AddUint96 adds a uint96 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint96() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint96,
		array: false,
	})
}

// AddUint104 adds a uint104 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint104() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint104,
		array: false,
	})
}

// AddUint112 adds a uint112 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint112() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint112,
		array: false,
	})
}

// AddUint120 adds a uint120 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint120() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint120,
		array: false,
	})
}

// AddUint128 adds a uint128 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint128() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint128,
		array: false,
	})
}

// AddUint136 adds a uint136 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint136() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint136,
		array: false,
	})
}

// AddUint144 adds a uint144 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint144() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint144,
		array: false,
	})
}

// AddUint152 adds a uint152 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint152() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint152,
		array: false,
	})
}

// AddUint160 adds a uint160 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint160() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint160,
		array: false,
	})
}

// AddUint168 adds a uint168 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint168() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint168,
		array: false,
	})
}

// AddUint176 adds a uint176 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint176() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint176,
		array: false,
	})
}

// AddUint184 adds a uint184 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint184() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint184,
		array: false,
	})
}

// AddUint192 adds a uint192 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint192() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint192,
		array: false,
	})
}

// AddUint200 adds a uint200 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint200() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint200,
		array: false,
	})
}

// AddUint208 adds a uint208 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint208() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint208,
		array: false,
	})
}

// AddUint216 adds a uint216 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint216() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint216,
		array: false,
	})
}

// AddUint224 adds a uint224 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint224() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint224,
		array: false,
	})
}

// AddUint232 adds a uint232 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint232() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint232,
		array: false,
	})
}

// AddUint240 adds a uint240 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint240() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint240,
		array: false,
	})
}

// AddUint248 adds a uint248 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint248() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint248,
		array: false,
	})
}

// AddUint256 adds a uint256 parameter to the selector.
func (selector *ContractFunctionSelector) AddUint256() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint256,
		array: false,
	})
}

// AddBytes adds a bytes parameter to the selector.
func (selector *ContractFunctionSelector) AddBytes() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBytes,
		array: false,
	})
}

// AddBytes32 adds a bytes32 parameter to the selector.
func (selector *ContractFunctionSelector) AddBytes32() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBytes32,
		array: false,
	})
}

// AddAddressArray adds an address[] parameter to the selector.
func (selector *ContractFunctionSelector) AddAddressArray() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aAddress,
		array: true,
	})
}

// AddBoolArray adds a bool[] parameter to the selector.
func (selector *ContractFunctionSelector) AddBoolArray() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBool,
		array: true,
	})
}

// AddStringArray adds a string[] parameter to the selector.
func (selector *ContractFunctionSelector) AddStringArray() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aString,
		array: true,
	})
}

// AddInt8Array adds an int8[] parameter to the selector.
func (selector *ContractFunctionSelector) AddInt8Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt8,
		array: true,
	})
}

// AddInt32Array adds an int32[] parameter to the selector.
func (selector *ContractFunctionSelector) AddInt32Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt32,
		array: true,
	})
}

// AddInt64Array adds an int64[] parameter to the selector.
func (selector *ContractFunctionSelector) AddInt64Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt64,
		array: true,
	})
}

// AddInt256Array adds an int256[] parameter to the selector.
func (selector *ContractFunctionSelector) AddInt256Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt256,
		array: true,
	})
}

// AddUint8Array adds a uint8[] parameter to the selector.
func (selector *ContractFunctionSelector) AddUint8Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint8,
		array: true,
	})
}

// AddUint32Array adds a uint32[] parameter to the selector.
func (selector *ContractFunctionSelector) AddUint32Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint32,
		array: true,
	})
}

// AddUint64Array adds a uint64[] parameter to the selector.
func (selector *ContractFunctionSelector) AddUint64Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint64,
		array: true,
	})
}

// AddUint256Array adds a uint256[] parameter to the selector.
func (selector *ContractFunctionSelector) AddUint256Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint256,
		array: true,
	})
}

// AddBytesArray adds a bytes[] parameter to the selector.
func (selector *ContractFunctionSelector) AddBytesArray() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBytes,
		array: true,
	})
}

// AddBytes32Array adds a bytes32[] parameter to the selector.
func (selector *ContractFunctionSelector) AddBytes32Array() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aBytes32,
		array: true,
	})
}

// String returns the string representation of the selector.
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
