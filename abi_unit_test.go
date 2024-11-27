//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbi(t *testing.T) {
	inputs, _ := NewType("tuple()")
	outputs, _ := NewType("tuple()")

	methodOutput := &Method{
		Name:    "abc",
		Inputs:  inputs,
		Outputs: outputs,
	}

	inputs, _ = NewType("tuple(address owner)")
	outputs, _ = NewType("tuple(uint256 balance)")
	balanceFunc := &Method{
		Name:    "balanceOf",
		Const:   true,
		Inputs:  inputs,
		Outputs: outputs,
	}

	errorInput, _ := NewType("tuple(address indexed a)")
	eventInput, _ := NewType("tuple(address indexed a)")

	cases := []struct {
		Input  string
		Output *ABI
	}{
		{
			Input: `[
				{
					"name": "abc",
					"type": "function"
				},
				{
					"name": "cde",
					"type": "event",
					"inputs": [
						{
							"indexed": true,
							"name": "a",
							"type": "address"
						}
					]
				},
				{
					"name": "def",
					"type": "error",
					"inputs": [
						{
							"indexed": true,
							"name": "a",
							"type": "address"
						}
					]
				},
				{
					"type": "function",
					"name": "balanceOf",
					"constant": true,
					"stateMutability": "view",
				 	"payable": false,
					"inputs": [
						{
					    	"type": "address",
					    	"name": "owner"
					   	}
					],
					"outputs": [
						{
					    	"type": "uint256",
					    	"name": "balance"
					   	}
					]
				}
			]`,
			Output: &ABI{
				Events: map[string]*Event{
					"cde": {
						Name:   "cde",
						Inputs: eventInput,
					},
				},
				Methods: map[string]*Method{
					"abc":       methodOutput,
					"balanceOf": balanceFunc,
				},
				MethodsBySignature: map[string]*Method{
					"abc()":              methodOutput,
					"balanceOf(address)": balanceFunc,
				},
				Errors: map[string]*Error{
					"def": {
						Name:   "def",
						Inputs: errorInput,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			abi, err := NewABI(c.Input)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(abi, c.Output) {
				t.Fatal("bad")
			}
		})
	}
}

func TestAbiInternalType(t *testing.T) {
	const abiStr = `[
        {
            "inputs": [
				{
					"components": [
						{
							"internalType": "address",
							"type": "address"
						},
						{
							"internalType": "uint256[4]",
							"type": "uint256[4]"
						}
					],
					"internalType": "struct X",
					"name": "newSet",
					"type": "tuple[]"
				},
                {
                    "internalType": "custom_address",
                    "name": "_to",
                    "type": "address"
                }
			],
			"outputs": [],
			"name": "transfer",
			"type": "function"
		}
	]`

	abi, err := NewABI(abiStr)
	require.NoError(t, err)

	typ := abi.GetMethod("transfer").Inputs
	require.Equal(t, typ.tuple[0].Elem.InternalType(), "struct X")
	require.Equal(t, typ.tuple[0].Elem.elem.tuple[0].Elem.InternalType(), "address")
	require.Equal(t, typ.tuple[0].Elem.elem.tuple[1].Elem.InternalType(), "uint256[4]")
	require.Equal(t, typ.tuple[1].Elem.InternalType(), "custom_address")
}

func TestAbiPolymorphism(t *testing.T) {
	// This ABI contains 2 "transfer" functions (polymorphism)
	const polymorphicABI = `[
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "_to",
                    "type": "address"
                },
                {
                    "internalType": "address",
                    "name": "_token",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "_amount",
                    "type": "uint256"
                }
            ],
            "name": "transfer",
            "outputs": [
                {
                    "internalType": "bool",
                    "name": "",
                    "type": "bool"
                }
            ],
            "stateMutability": "nonpayable",
            "type": "function"
        },
		{
            "inputs": [
                {
                    "internalType": "address",
                    "name": "_to",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "_amount",
                    "type": "uint256"
                }
            ],
            "name": "transfer",
            "outputs": [
                {
                    "internalType": "bool",
                    "name": "",
                    "type": "bool"
                }
            ],
            "stateMutability": "nonpayable",
            "type": "function"
        }
    ]`

	abi, err := NewABI(polymorphicABI)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, abi.Methods, 2)
	assert.Equal(t, abi.GetMethod("transfer").Sig(), "transfer(address,address,uint256)")
	assert.Equal(t, abi.GetMethod("transfer0").Sig(), "transfer(address,uint256)")
	assert.NotEmpty(t, abi.GetMethodBySignature("transfer(address,address,uint256)"))
	assert.NotEmpty(t, abi.GetMethodBySignature("transfer(address,uint256)"))
}

func TestAbiHumanReadable(t *testing.T) {
	cases := []string{
		"constructor(string symbol, string name)",
		"function transferFrom(address from, address to, uint256 value)",
		"function balanceOf(address owner) view returns (uint256 balance)",
		"function balanceOf() view returns ()",
		"event Transfer(address indexed from, address indexed to, address value)",
		"error InsufficientBalance(address owner, uint256 balance)",
		"function addPerson(tuple(string name, uint16 age) person)",
		"function addPeople(tuple(string name, uint16 age)[] person)",
		"function getPerson(uint256 id) view returns (tuple(string name, uint16 age))",
		"event PersonAdded(uint256 indexed id, tuple(string name, uint16 age) person)",
	}
	vv, err := NewABIFromList(cases)
	assert.NoError(t, err)

	// make it nil to not compare it and avoid writing each method twice for the test
	vv.MethodsBySignature = nil

	constructorInputs, _ := NewType("tuple(string symbol, string name)")
	transferFromInputs, _ := NewType("tuple(address from, address to, uint256 value)")
	transferFromOutputs, _ := NewType("tuple()")
	balanceOfInputs, _ := NewType("tuple(address owner)")
	balanceOfOutputs, _ := NewType("tuple(uint256 balance)")
	balanceOf0Inputs, _ := NewType("tuple()")
	balanceOf0Outputs, _ := NewType("tuple()")
	addPersonInputs, _ := NewType("tuple(tuple(string name, uint16 age) person)")
	addPersonOutputs, _ := NewType("tuple()")
	addPeopleInputs, _ := NewType("tuple(tuple(string name, uint16 age)[] person)")
	addPeopleOutputs, _ := NewType("tuple()")
	getPersonInputs, _ := NewType("tuple(uint256 id)")
	getPersonOutputs, _ := NewType("tuple(tuple(string name, uint16 age))")
	transferEventInputs, _ := NewType("tuple(address indexed from, address indexed to, address value)")
	personAddedEventInputs, _ := NewType("tuple(uint256 indexed id, tuple(string name, uint16 age) person)")
	errorInputs, _ := NewType("tuple(address owner, uint256 balance)")

	expect := &ABI{
		Constructor: &Method{
			Inputs: constructorInputs,
		},
		Methods: map[string]*Method{
			"transferFrom": {
				Name:    "transferFrom",
				Inputs:  transferFromInputs,
				Outputs: transferFromOutputs,
			},
			"balanceOf": {
				Name:    "balanceOf",
				Inputs:  balanceOfInputs,
				Outputs: balanceOfOutputs,
			},
			"balanceOf0": {
				Name:    "balanceOf",
				Inputs:  balanceOf0Inputs,
				Outputs: balanceOf0Outputs,
			},
			"addPerson": {
				Name:    "addPerson",
				Inputs:  addPersonInputs,
				Outputs: addPersonOutputs,
			},
			"addPeople": {
				Name:    "addPeople",
				Inputs:  addPeopleInputs,
				Outputs: addPeopleOutputs,
			},
			"getPerson": {
				Name:    "getPerson",
				Inputs:  getPersonInputs,
				Outputs: getPersonOutputs,
			},
		},
		Events: map[string]*Event{
			"Transfer": {
				Name:   "Transfer",
				Inputs: transferEventInputs,
			},
			"PersonAdded": {
				Name:   "PersonAdded",
				Inputs: personAddedEventInputs,
			},
		},
		Errors: map[string]*Error{
			"InsufficientBalance": {
				Name:   "InsufficientBalance",
				Inputs: errorInputs,
			},
		},
	}
	assert.Equal(t, expect, vv)
}

func TestAbiParseMethodSignature(t *testing.T) {
	cases := []struct {
		signature string
		name      string
		input     string
		output    string
	}{
		{
			// both input and output
			signature: "function approve(address to) returns (address)",
			name:      "approve",
			input:     "tuple(address)",
			output:    "tuple(address)",
		},
		{
			// no input
			signature: "function approve() returns (address)",
			name:      "approve",
			input:     "tuple()",
			output:    "tuple(address)",
		},
		{
			// no output
			signature: "function approve(address)",
			name:      "approve",
			input:     "tuple(address)",
			output:    "tuple()",
		},
		{
			// multiline
			signature: `function a(
				uint256 b,
				address[] c
			)
				returns
				(
				uint256[] d
			)`,
			name:   "a",
			input:  "tuple(uint256,address[])",
			output: "tuple(uint256[])",
		},
	}

	for _, c := range cases {
		name, input, output, err := parseMethodSignature(c.signature)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, name, c.name)

		if input != nil {
			assert.Equal(t, c.input, input.String())
		} else {
			assert.Equal(t, c.input, "")
		}

		if input != nil {
			assert.Equal(t, c.output, output.String())
		} else {
			assert.Equal(t, c.output, "")
		}
	}
}
