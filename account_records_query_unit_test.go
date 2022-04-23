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

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountRecordQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	recordQuery := NewAccountRecordsQuery().
		SetAccountID(accountID)

	err = recordQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountRecordQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	recordQuery := NewAccountRecordsQuery().
		SetAccountID(accountID)

	err = recordQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitMockAccountRecordsQuery(t *testing.T) {
	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_CryptoGetAccountRecords{
				CryptoGetAccountRecords: &services.CryptoGetAccountRecordsResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
		&services.Response{
			Response: &services.Response_CryptoGetAccountRecords{
				CryptoGetAccountRecords: &services.CryptoGetAccountRecordsResponse{
					Header:    &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 1},
					AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 1800}},
					Records: []*services.TransactionRecord{
						{
							TransactionHash:    []byte{1},
							ConsensusTimestamp: &services.Timestamp{Nanos: 12313123, Seconds: 2313},
							TransactionID: &services.TransactionID{
								TransactionValidStart: &services.Timestamp{Nanos: 12313123, Seconds: 2313},
								AccountID:             &services.AccountID{Account: &services.AccountID_AccountNum{AccountNum: 1800}},
								Scheduled:             false,
								Nonce:                 0,
							},
							Memo:           "",
							TransactionFee: 0,
						},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	recordsQuery, err := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAccountID(AccountID{Account: 1800}).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, len(recordsQuery), 1)
	require.Equal(t, recordsQuery[0].TransactionID.AccountID.Account, uint64(1800))
}

func TestUnitAccountRecordsQueryGet(t *testing.T) {
	spenderAccountID1 := AccountID{Account: 7}

	balance := NewAccountRecordsQuery().
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

func TestUnitAccountRecordsQuerySetNothing(t *testing.T) {
	balance := NewAccountRecordsQuery()

	balance.GetAccountID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}
