//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitContractCallQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	contractCall := NewContractCallQuery().
		SetContractID(contractID)

	err = contractCall.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitContractCallQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	contractCall := NewContractCallQuery().
		SetContractID(contractID)

	err = contractCall.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitMockContractCallQuery(t *testing.T) {
	t.Parallel()

	message := "getMessage"
	params := NewContractFunctionParameters()
	params._Build(&message)

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_ContractCallLocal{
				ContractCallLocal: &services.ContractCallLocalResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ContractCallLocal{
				ContractCallLocal: &services.ContractCallLocalResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					FunctionResult: &services.ContractFunctionResult{
						ContractID:         &services.ContractID{Contract: &services.ContractID_ContractNum{ContractNum: 123}},
						GasUsed:            75000,
						ContractCallResult: params._Build(&message),
						SignerNonce:        wrapperspb.Int64(0),
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	result, err := NewContractCallQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetContractID(ContractID{Contract: 123}).
		SetQueryPayment(NewHbar(1)).
		SetGas(100000).
		SetFunction(message, nil).
		SetMaxQueryPayment(NewHbar(5)).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, bytes.Compare(result.ContractCallResult, params._Build(&message)), 0)
	require.Equal(t, result.GasUsed, uint64(75000))
	require.Equal(t, result.ContractID.Contract, uint64(123))
}

func TestUnitContractCallQueryGet(t *testing.T) {
	t.Parallel()

	spenderContractID := ContractID{Contract: 7}

	balance := NewContractCallQuery().
		SetContractID(spenderContractID).
		SetQueryPayment(NewHbar(2)).
		SetGas(100000).
		SetFunction("getMessage", nil).
		SetFunctionParameters([]byte{}).
		SetMaxQueryPayment(NewHbar(10)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetContractID()
	balance.GetFunctionParameters()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitContractCallQuerySetNothing(t *testing.T) {
	t.Parallel()

	balance := NewContractCallQuery()

	balance.GetContractID()
	balance.GetFunctionParameters()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}
