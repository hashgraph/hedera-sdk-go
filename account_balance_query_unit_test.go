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

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountBalanceQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	balanceQuery := NewAccountBalanceQuery().
		SetAccountID(accountID)

	err = balanceQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountBalanceQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	balanceQuery := NewAccountBalanceQuery().
		SetAccountID(accountID)

	err = balanceQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountBalanceQueryGet(t *testing.T) {
	t.Parallel()

	spenderAccountID1 := AccountID{Account: 7}

	balance := NewAccountBalanceQuery().
		SetAccountID(spenderAccountID1).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetAccountID()
	balance.GetNodeAccountIDs()
	balance.GetPaymentTransactionID()
}

func TestUnitAccountBalanceQuerySetNothing(t *testing.T) {
	t.Parallel()

	balance := NewAccountBalanceQuery()

	balance.GetAccountID()
	balance.GetNodeAccountIDs()
	balance.GetPaymentTransactionID()
}

func TestUnitAccountBalanceQueryCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	contract := ContractID{Contract: 3, checksum: &checksum}
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewAccountBalanceQuery().
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetAccountID(account).
		SetContractID(contract).
		SetNodeAccountIDs(nodeAccountID).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(23)).
		SetQueryPayment(NewHbar(3))

	err = query._ValidateNetworkOnIDs(client)

	require.NoError(t, err)
	query.GetNodeAccountIDs()
	query.GetMaxBackoff()
	query.GetMinBackoff()
	query._GetLogID()
	query.GetAccountID()
	query.GetContractID()

	_AccountBalanceFromProtobuf(nil)
	bal := AccountBalance{Hbars: NewHbar(2)}
	bal._ToProtobuf()
}

func TestUnitAccountBalanceQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{
		{
			&services.Response{
				Response: &services.Response_CryptogetAccountBalance{
					CryptogetAccountBalance: &services.CryptoGetAccountBalanceResponse{
						Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 0},
						AccountID: &services.AccountID{ShardNum: 0, RealmNum: 0, Account: &services.AccountID_AccountNum{
							AccountNum: 1800,
						}},
						Balance: 2000,
					},
				},
			},
			&services.Response{
				Response: &services.Response_CryptogetAccountBalance{
					CryptogetAccountBalance: &services.CryptoGetAccountBalanceResponse{
						Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 0},
						AccountID: &services.AccountID{ShardNum: 0, RealmNum: 0, Account: &services.AccountID_AccountNum{
							AccountNum: 1800,
						}},
						Balance: 2000,
					},
				},
			},
		},
	}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewAccountBalanceQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAccountID(AccountID{Account: 1800}).
		SetContractID(ContractID{Contract: 3})

	query.GetCost(client)
	_, err := query.Execute(client)
	require.NoError(t, err)
}
