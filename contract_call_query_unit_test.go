//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"bytes"
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"

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

	err = contractCall._ValidateNetworkOnIDs(client)
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

	err = contractCall._ValidateNetworkOnIDs(client)
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
