//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationContractIDCanPopulateAccountNumber(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	resp, err := NewFileCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := *receipt.FileID
	require.NotNil(t, fileID)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetGas(100000).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hiero")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hiero-sdk-go::TestContractInfoQuery_Execute").
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	require.NotNil(t, receipt.ContractID)

	contractID := *receipt.ContractID
	info, err := NewContractInfoQuery().SetContractID(contractID).Execute(env.Client)
	require.NoError(t, err)
	idMirror, err := ContractIDFromEvmAddress(0, 0, info.ContractAccountID)
	require.NoError(t, err)
	time.Sleep(5 * time.Second)
	err = idMirror.PopulateContract(env.Client)
	require.NoError(t, err)
	require.Equal(t, contractID.Contract, idMirror.Contract)
}
