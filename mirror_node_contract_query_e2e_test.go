//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	SMART_CONTRACT_BYTECODE = "6080604052348015600e575f80fd5b50335f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506104a38061005b5f395ff3fe608060405260043610610033575f3560e01c8063607a4427146100375780637065cb4814610053578063893d20e81461007b575b5f80fd5b610051600480360381019061004c919061033c565b6100a5565b005b34801561005e575f80fd5b50610079600480360381019061007491906103a2565b610215565b005b348015610086575f80fd5b5061008f6102b7565b60405161009c91906103dc565b60405180910390f35b3373ffffffffffffffffffffffffffffffffffffffff165f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146100fb575f80fd5b805f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600181908060018154018082558091505060019003905f5260205f20015f9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505f8173ffffffffffffffffffffffffffffffffffffffff166108fc3490811502906040515f60405180830381858888f19350505050905080610211576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102089061044f565b60405180910390fd5b5050565b805f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600181908060018154018082558091505060019003905f5260205f20015f9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61030b826102e2565b9050919050565b61031b81610301565b8114610325575f80fd5b50565b5f8135905061033681610312565b92915050565b5f60208284031215610351576103506102de565b5b5f61035e84828501610328565b91505092915050565b5f610371826102e2565b9050919050565b61038181610367565b811461038b575f80fd5b50565b5f8135905061039c81610378565b92915050565b5f602082840312156103b7576103b66102de565b5b5f6103c48482850161038e565b91505092915050565b6103d681610367565b82525050565b5f6020820190506103ef5f8301846103cd565b92915050565b5f82825260208201905092915050565b7f5472616e73666572206661696c656400000000000000000000000000000000005f82015250565b5f610439600f836103f5565b915061044482610405565b602082019050919050565b5f6020820190508181035f8301526104668161042d565b905091905056fea26469706673582212206c46ddb2acdbcc4290e15be83eb90cd0b2ce5bd82b9bfe58a0709c5aec96305564736f6c634300081a0033"
	ADDRESS                 = "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"
)

func TestMirrorNodeContractQueryCanSimulateTransaction(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey).
		SetContents([]byte(SMART_CONTRACT_BYTECODE)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := receipt.FileID

	response, err = NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetBytecodeFileID(*fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	// Wait for mirror node to import data
	time.Sleep(2 * time.Second)

	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.NoError(t, err)

	result, err := NewContractCallQuery().
		SetContractID(*contractID).
		SetGas(gas).
		SetFunction("getOwner", nil).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	simulationResult, err := NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("%x", result.GetAddress(0)), simulationResult[26:])

	param, err := NewContractFunctionParameters().AddAddress(ADDRESS)
	require.NoError(t, err)

	gas, err = NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetFunction("addOwner", param).
		Execute(env.Client)
	require.NoError(t, err)

	response, err = NewContractExecuteTransaction().
		SetContractID(*contractID).
		SetGas(gas).
		SetFunction("addOwner", param).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetFunction("addOwner", param).
		Execute(env.Client)
	require.NoError(t, err)
}

func TestMirrorNodeContractQueryReturnsDefaultGasWhenContractIsNotDeployed(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	contractID := ContractID{Contract: 1231456}
	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(contractID).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.NoError(t, err)
	require.Equal(t, uint64(22892), gas)
}

func TestMirrorNodeContractQueryFailWhenGasLimitIsLow(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey).
		SetContents([]byte(SMART_CONTRACT_BYTECODE)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := receipt.FileID

	response, err = NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetBytecodeFileID(*fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	// Wait for mirror node to import data
	time.Sleep(2 * time.Second)

	_, err = NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetGasLimit(100).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.ErrorContains(t, err, "received non-200 response from Mirror Node")

	_, err = NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetGasLimit(100).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.ErrorContains(t, err, "received non-200 response from Mirror Node")
}

func TestMirrorNodeContractQueryFailWhenSenderIsNotSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey).
		SetContents([]byte(SMART_CONTRACT_BYTECODE)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := receipt.FileID

	response, err = NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetBytecodeFileID(*fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	// Wait for mirror node to import data
	time.Sleep(2 * time.Second)
	param, err := NewContractFunctionParameters().AddAddress(ADDRESS)

	_, err = NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetFunction("addOwnerAndTransfer", param).
		Execute(env.Client)
	require.ErrorContains(t, err, "received non-200 response from Mirror Node")

	_, err = NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetFunction("addOwnerAndTransfer", param).
		Execute(env.Client)
	require.ErrorContains(t, err, "received non-200 response from Mirror Node")
}
func TestMirrorNodeContractQueryCanSimulateWithSenderSet(t *testing.T) {
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	response, err := NewFileCreateTransaction().
		SetKeys(env.OperatorKey).
		SetContents([]byte(SMART_CONTRACT_BYTECODE)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := receipt.FileID

	response, err = NewContractCreateTransaction().
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetBytecodeFileID(*fileID).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	receiverId, _, err := createAccount(&env, func(transaction *AccountCreateTransaction) {
		pk, _ := PrivateKeyGenerateEd25519()
		transaction.SetKey(pk)
	})
	require.NoError(t, err)
	receiverEvmAddress := receiverId.ToSolidityAddress()

	// Wait for mirror node to import data
	time.Sleep(2 * time.Second)
	param, err := NewContractFunctionParameters().AddAddress(receiverEvmAddress)

	owner, err := NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetFunction("getOwner", nil).
		Execute(env.Client)
	require.NoError(t, err)

	gas, err := NewMirrorNodeContractEstimateGasQuery().
		SetContractID(*contractID).
		SetSenderEvmAddress(owner[26:]).
		SetFunction("addOwnerAndTransfer", param).
		SetValue(123).
		SetGasLimit(1_000_000).
		Execute(env.Client)
	require.NoError(t, err)

	resp, err := NewContractExecuteTransaction().
		SetContractID(*contractID).
		SetGas(gas).
		SetPayableAmount(NewHbar(1)).
		SetFunction("addOwnerAndTransfer", param).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = NewMirrorNodeContractCallQuery().
		SetContractID(*contractID).
		SetSenderEvmAddress(owner[26:]).
		SetFunction("addOwnerAndTransfer", param).
		SetValue(123).
		SetGasLimit(1_000_000).
		Execute(env.Client)
	require.NoError(t, err)
}
