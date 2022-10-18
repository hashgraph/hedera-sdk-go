//go:build all || unit
// +build all unit

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
	"testing"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountStakersQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	stackersQuery := NewAccountStakersQuery().
		SetAccountID(accountID)

	err = stackersQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountStakersQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	stackersQuery := NewAccountStakersQuery().
		SetAccountID(accountID)

	err = stackersQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountStakersQueryGet(t *testing.T) {
	spenderAccountID1 := AccountID{Account: 7}

	balance := NewAccountStakersQuery().
		SetAccountID(spenderAccountID1).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(10)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetAccountID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitAccountStakersQuerySetNothing(t *testing.T) {
	balance := NewAccountStakersQuery()

	balance.GetAccountID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitAccountStakersQueryCoverage(t *testing.T) {
	checksum := "dmqui"
	grpc := time.Second * 3
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)

	query := NewAccountStakersQuery().
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetAccountID(account).
		SetNodeAccountIDs(nodeAccountID).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(23)).
		SetQueryPayment(NewHbar(3)).
		SetGrpcDeadline(&grpc)

	err := query._ValidateNetworkOnIDs(client)

	require.NoError(t, err)
	query.GetNodeAccountIDs()
	query.GetMaxBackoff()
	query.GetMinBackoff()
	query._GetLogID()
	query.GetAccountID()
}

func TestUnitAccountStakersQueryQueryMock(t *testing.T) {
	responses := [][]interface{}{
		{
			&services.Response{
				Response: &services.Response_CryptoGetProxyStakers{
					CryptoGetProxyStakers: &services.CryptoGetStakersResponse{
						Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 0},
					},
				},
			},
			&services.Response{
				Response: &services.Response_CryptoGetProxyStakers{
					CryptoGetProxyStakers: &services.CryptoGetStakersResponse{
						Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 0},
					},
				},
			},
			&services.Response{
				Response: &services.Response_CryptoGetProxyStakers{
					CryptoGetProxyStakers: &services.CryptoGetStakersResponse{
						Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 0},
						Stakers: &services.AllProxyStakers{
							AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 1800}},
						},
					},
				},
			},
		},
	}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewAccountStakersQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAccountID(AccountID{Account: 1800})

	query.GetCost(client)
	_, err := query.Execute(client)
	require.NoError(t, err)
}
