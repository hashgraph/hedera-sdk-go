//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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
	"math/big"
	"reflect"
	"testing"
)

func TestEncodeStruct(t *testing.T) {
	typ, _ := NewType("tuple(address aa, uint256 b)")

	type Obj struct {
		A Address `abi:"aa"`
		B *big.Int
	}
	obj := Obj{
		A: Address{0x1},
		B: big.NewInt(1),
	}

	encoded, err := typ.Encode(&obj)
	if err != nil {
		t.Fatal(err)
	}

	var obj2 Obj
	if err := typ.DecodeStruct(encoded, &obj2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(obj, obj2) {
		t.Fatal("bad")
	}
}

func TestEncodeStructCamcelCase(t *testing.T) {
	typ, _ := NewType("tuple(address aA, uint256 b)")

	type Obj struct {
		A Address `abi:"aA"`
		B *big.Int
	}
	obj := Obj{
		A: Address{0x1},
		B: big.NewInt(1),
	}

	encoded, err := typ.Encode(&obj)
	if err != nil {
		t.Fatal(err)
	}

	var obj2 Obj
	if err := typ.DecodeStruct(encoded, &obj2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(obj, obj2) {
		t.Fatal("bad")
	}
}
