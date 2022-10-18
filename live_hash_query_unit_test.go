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

func TestUnitLiveHashQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	liveHashQuery := NewLiveHashQuery().
		SetAccountID(accountID)

	err = liveHashQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitLiveHashQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	liveHashQuery := NewLiveHashQuery().
		SetAccountID(accountID)

	err = liveHashQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitLiveHashQueryGet(t *testing.T) {
	accountID := AccountID{Account: 7}

	balance := NewLiveHashQuery().
		SetAccountID(accountID).
		SetHash([]byte{}).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetAccountID()
	balance.GetGetHash()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitLiveHashQuerySetNothing(t *testing.T) {
	balance := NewLiveHashQuery()

	balance.GetAccountID()
	balance.GetGetHash()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitLiveHashQueryCoverage(t *testing.T) {
	checksum := "dmqui"
	grpc := time.Second * 3
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)

	query := NewLiveHashQuery().
		SetAccountID(account).
		SetHash([]byte{1}).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
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
	query.GetGetHash()
	query.GetQueryPayment()
	query.GetMaxQueryPayment()
}

func TestUnitLiveHashQueryMock(t *testing.T) {
	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_CryptoGetLiveHash{
				CryptoGetLiveHash: &services.CryptoGetLiveHashResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_CryptoGetLiveHash{
				CryptoGetLiveHash: &services.CryptoGetLiveHashResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_CryptoGetLiveHash{
				CryptoGetLiveHash: &services.CryptoGetLiveHashResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					LiveHash: &services.LiveHash{
						Hash:      []byte{1},
						AccountId: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 3}},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewLiveHashQuery().
		SetAccountID(AccountID{Account: 3}).
		SetHash([]byte{1}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	_, err := query.GetCost(client)
	require.NoError(t, err)

	_, err = query.Execute(client)
	require.NoError(t, err)
}
