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

func (selector *ContractFunctionSelector) AddInt16() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt16,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt24() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt24,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt32() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt32,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt40() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt40,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt48() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt48,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt56() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt56,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt64() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt64,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt72() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt72,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt80() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt80,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt88() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt88,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt96() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt96,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt104() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt104,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt112() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt112,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt120() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt120,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt128() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt128,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt136() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt136,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt144() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt144,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt152() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt152,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt160() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt160,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddInt168() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt168,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt176() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt176,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt184() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt184,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt192() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt192,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt200() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt200,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt208() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt208,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt216() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt216,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt224() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt224,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt232() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt232,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt240() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt240,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddInt248() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aInt248,
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

func (selector *ContractFunctionSelector) AddUint16() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint16,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint24() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint24,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint32() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint32,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint40() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint40,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint48() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint48,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint56() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint56,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint64() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint64,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint72() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint72,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint80() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint80,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint88() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint88,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint96() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint96,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint104() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint104,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint112() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint112,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint120() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint120,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint128() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint128,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint136() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint136,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint144() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint144,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint152() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint152,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint160() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint160,
		array: false,
	})
}
func (selector *ContractFunctionSelector) AddUint168() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint168,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint176() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint176,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint184() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint184,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint192() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint192,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint200() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint200,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint208() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint208,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint216() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint216,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint224() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint224,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint232() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint232,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint240() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint240,
		array: false,
	})
}

func (selector *ContractFunctionSelector) AddUint248() *ContractFunctionSelector {
	return selector._AddParam(_Solidity{
		ty:    aUint248,
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
