//go:build all || e2e
// +build all e2e

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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationContractCreateFlowCanExecute(t *testing.T) {
	testContractByteCode := []byte(`608060405234801561001057600080fd5b506040516104d73803806104d78339818101604052602081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825164010000000081118282018810171561008257600080fd5b82525081516020918201929091019080838360005b838110156100af578181015183820152602001610097565b50505050905090810190601f1680156100dc5780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b0319163317905550805161010890600190602084019061010f565b50506101aa565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015057805160ff191683800117855561017d565b8280016001018555821561017d579182015b8281111561017d578251825591602001919060010190610162565b5061018992915061018d565b5090565b6101a791905b808211156101895760008155600101610193565b90565b61031e806101b96000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063368b87721461004657806341c0e1b5146100ee578063ce6d41de146100f6575b600080fd5b6100ec6004803603602081101561005c57600080fd5b81019060208101813564010000000081111561007757600080fd5b82018360208201111561008957600080fd5b803590602001918460018302840111640100000000831117156100ab57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610173945050505050565b005b6100ec6101a2565b6100fe6101ba565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610138578181015183820152602001610120565b50505050905090810190601f1680156101655780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000546001600160a01b0316331461018a5761019f565b805161019d906001906020840190610250565b505b50565b6000546001600160a01b03163314156101b85733ff5b565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156102455780601f1061021a57610100808354040283529160200191610245565b820191906000526020600020905b81548152906001019060200180831161022857829003601f168201915b505050505090505b90565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061029157805160ff19168380011785556102be565b828001600101855582156102be579182015b828111156102be5782518255916020019190600101906102a3565b506102ca9291506102ce565b5090565b61024d91905b808211156102ca57600081556001016102d456fea264697066735822122084964d4c3f6bc912a9d20e14e449721012d625aa3c8a12de41ae5519752fc89064736f6c63430006000033`)

	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewContractCreateFlow().
		SetBytecode(testContractByteCode).
		SetAdminKey(env.OperatorKey).
		SetGas(200000).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetContractMemo("[e2e::ContractCreateFlow]").
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	require.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	resp, err = NewContractExecuteTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetGas(200000).
		SetFunction("setMessage", NewContractFunctionParameters().AddString("new message")).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(AccountID{Account: 3}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationContractCreateFlowSmallContractCanExecute(t *testing.T) {
	t.Skip("Possible bug in services is causing this test to behave flaky")

	testContractByteCode := []byte(`6080604052348015600f57600080fd5b50604051601a90603b565b604051809103906000f0801580156035573d6000803e3d6000fd5b50506047565b605c8061009483390190565b603f806100556000396000f3fe6080604052600080fdfea2646970667358221220a20122cbad3457fedcc0600363d6e895f17048f5caa4afdab9e655123737567d64736f6c634300081200336080604052348015600f57600080fd5b50603f80601d6000396000f3fe6080604052600080fdfea264697066735822122053dfd8835e3dc6fedfb8b4806460b9b7163f8a7248bac510c6d6808d9da9d6d364736f6c63430008120033`)

	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewContractCreateFlow().
		SetBytecode(testContractByteCode).
		SetAdminKey(env.OperatorKey).
		SetGas(100000).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	require.NotNil(t, receipt.ContractID)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationContractCreateFlowGettersAndSetters(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	testContractByteCode := "608060405234801561001057600080fd5b506040516104d73803806104d78339818101604052602081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825164010000000081118282018810171561008257600080fd5b82525081516020918201929091019080838360005b838110156100af578181015183820152602001610097565b50505050905090810190601f1680156100dc5780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b0319163317905550805161010890600190602084019061010f565b50506101aa565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015057805160ff191683800117855561017d565b8280016001018555821561017d579182015b8281111561017d578251825591602001919060010190610162565b5061018992915061018d565b5090565b6101a791905b808211156101895760008155600101610193565b90565b61031e806101b96000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063368b87721461004657806341c0e1b5146100ee578063ce6d41de146100f6575b600080fd5b6100ec6004803603602081101561005c57600080fd5b81019060208101813564010000000081111561007757600080fd5b82018360208201111561008957600080fd5b803590602001918460018302840111640100000000831117156100ab57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610173945050505050565b005b6100ec6101a2565b6100fe6101ba565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610138578181015183820152602001610120565b50505050905090810190601f1680156101655780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000546001600160a01b0316331461018a5761019f565b805161019d906001906020840190610250565b505b50565b6000546001600160a01b03163314156101b85733ff5b565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156102455780601f1061021a57610100808354040283529160200191610245565b820191906000526020600020905b81548152906001019060200180831161022857829003601f168201915b505050505090505b90565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061029157805160ff19168380011785556102be565b828001600101855582156102be579182015b828111156102be5782518255916020019190600101906102a3565b506102ca9291506102ce565b5090565b61024d91905b808211156102ca57600081556001016102d456fea264697066735822122084964d4c3f6bc912a9d20e14e449721012d625aa3c8a12de41ae5519752fc89064736f6c63430006000033"
	testOperatorKey := env.OperatorKey
	testGas := int64(100000)
	testMemo := "[e2e::ContractCreateFlow]"
	testInitialBalance := HbarFromTinybar(1000000000000000000)
	autoRenewPeriod := 131500 * time.Minute
	testProxyAccountID := AccountID{Account: 3}
	testAutoTokenAssociations := int32(10)
	testMaxChuncks := uint64(4)
	testNodeAccountID := AccountID{Account: 55}

	constructorFuncName := "constructor"
	constructorParams := ContractFunctionParameters{
		function: ContractFunctionSelector{
			function:   &constructorFuncName,
			params:     "wooow",
			paramTypes: []_Solidity{}},
		arguments: []Argument{Argument{
			value:   []byte("hello from hedera"),
			dynamic: false,
		}}}

	createChain := NewContractCreateFlow().
		SetBytecodeWithString(testContractByteCode).
		SetAdminKey(env.OperatorKey).
		SetGas(100000).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetContractMemo("[e2e::ContractCreateFlow]").
		SetInitialBalance(testInitialBalance).
		SetProxyAccountID(testProxyAccountID).
		SetConstructorParameters(&constructorParams).
		SetAutoRenewAccountID(testProxyAccountID).
		SetMaxAutomaticTokenAssociations(testAutoTokenAssociations).
		SetMaxChunks(testMaxChuncks).
		SetNodeAccountIDs([]AccountID{testNodeAccountID, testProxyAccountID})

	require.Equal(t, createChain.GetBytecode(), testContractByteCode)
	require.Equal(t, createChain.GetAdminKey(), testOperatorKey)
	require.Equal(t, createChain.GetGas(), testGas)
	require.Equal(t, createChain.GetContractMemo(), testMemo)
	require.Equal(t, createChain.GetInitialBalance(), testInitialBalance)
	require.Equal(t, createChain.GetAutoRenewPeriod(), autoRenewPeriod)
	require.Equal(t, createChain.GetProxyAccountID(), testProxyAccountID)
	require.Equal(t, createChain.GetAutoRenewAccountID(), testProxyAccountID)
	require.Equal(t, createChain.GetMaxAutomaticTokenAssociations(), testAutoTokenAssociations)
	require.Equal(t, createChain.GetMaxChunks(), testMaxChuncks)
	require.Equal(t, createChain.GetNodeAccountIDs(), []AccountID{testNodeAccountID, testProxyAccountID})

}
