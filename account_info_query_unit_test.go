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
	"testing"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/require"
)

func TestUnitAccountInfoQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	infoQuery := NewAccountInfoQuery().
		SetAccountID(accountID)

	err = infoQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}
func TestAccountInfoQuery_Get(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	deadline := time.Duration(time.Minute)
	accountId := AccountID{Account: 3, checksum: &checksum}
	transactionID := TransactionIDGenerate(accountId)
	query := NewAccountInfoQuery().
		SetAccountID(accountId).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(10)).
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetMaxRetry(5).
		SetMaxBackoff(10 * time.Second).
		SetMinBackoff(1 * time.Second).
		SetPaymentTransactionID(transactionID).
		SetGrpcDeadline(&deadline)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	err = query._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
	require.Equal(t, accountId, query.GetAccountID())
	require.Equal(t, NewHbar(2), query.GetQueryPayment())
	require.Equal(t, NewHbar(10), query.GetMaxQueryPayment())
	require.Equal(t, []AccountID{{Account: 3}, {Account: 4}}, query.GetNodeAccountIDs())
	require.Equal(t, 5, query.GetMaxRetryCount())
	require.Equal(t, 10*time.Second, query.GetMaxBackoff())
	require.Equal(t, 1*time.Second, query.GetMinBackoff())
	require.Equal(t, transactionID, query.GetPaymentTransactionID())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitAccountInfoQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	infoQuery := NewAccountInfoQuery().
		SetAccountID(accountID)

	err = infoQuery._ValidateNetworkOnIDs(client)
	require.Error(t, err)
	if err != nil {
		require.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountInfoQuerySetNothing(t *testing.T) {
	t.Parallel()

	balance := NewAccountInfoQuery()

	require.Equal(t, AccountID{}, balance.GetAccountID())
	require.Equal(t, []AccountID{}, balance.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, balance.GetMinBackoff())
	require.Equal(t, 8*time.Second, balance.GetMaxBackoff())
	require.Equal(t, 10, balance.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, balance.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, balance.GetQueryPayment())
	require.Equal(t, Hbar{}, balance.GetMaxQueryPayment())
}

func Test_AccountInfoQueryMapStatusError(t *testing.T) {
	t.Parallel()

	response := services.Response{
		Response: &services.Response_CryptoGetInfo{
			CryptoGetInfo: &services.CryptoGetInfoResponse{
				Header: &services.ResponseHeader{
					NodeTransactionPrecheckCode: services.ResponseCodeEnum(StatusInvalidAccountID),
					ResponseType:                services.ResponseType_COST_ANSWER,
				},
			},
		},
	}

	actualError := _AccountInfoQueryMapStatusError(nil, &response)

	expectedError := ErrHederaPreCheckStatus{
		Status: StatusInvalidAccountID,
	}

	require.Equal(t, expectedError, actualError)
}

func TestUnitAccountInfoQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_CryptoGetInfo{
				CryptoGetInfo: &services.CryptoGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 1234}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, HbarFromTinybar(2), cost)
}
