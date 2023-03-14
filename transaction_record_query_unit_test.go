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
	"fmt"
	"testing"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/require"
)

func TestUnitTransactionRecordQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	transactionID := TransactionIDGenerate(accountID)
	require.NoError(t, err)

	recordQuery := NewTransactionRecordQuery().
		SetTransactionID(transactionID)

	err = recordQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransactionRecordQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	transactionID := TransactionIDGenerate(accountID)
	require.NoError(t, err)

	recordQuery := NewTransactionRecordQuery().
		SetTransactionID(transactionID)

	err = recordQuery._ValidateNetworkOnIDs(client)
	require.Error(t, err)
	if err != nil {
		require.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTransactionRecordQueryGet(t *testing.T) {
	txID := TransactionIDGenerate(AccountID{Account: 7})
	deadline := time.Duration(time.Minute)
	accountId := AccountID{Account: 123}
	transactionID := TransactionIDGenerate(accountId)
	query := NewTransactionRecordQuery().
		SetTransactionID(txID).
		SetIncludeDuplicates(true).
		SetIncludeChildren(true).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}}).
		SetMaxRetry(3).
		SetMinBackoff(300 * time.Millisecond).
		SetMaxBackoff(10 * time.Second).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(500)).
		SetGrpcDeadline(&deadline)

	require.Equal(t, txID, query.GetTransactionID())
	require.True(t, query.GetIncludeChildren())
	require.True(t, query.GetIncludeDuplicates())
	require.Equal(t, []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}, query.GetNodeAccountIDs())
	require.Equal(t, 300*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 10*time.Second, query.GetMaxBackoff())
	require.Equal(t, 3, query.GetMaxRetryCount())
	require.Equal(t, transactionID, query.GetPaymentTransactionID())
	require.Equal(t, HbarFromTinybar(25), query.GetQueryPayment())
	require.Equal(t, NewHbar(500), query.GetMaxQueryPayment())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
	require.Equal(t, fmt.Sprintf("TransactionRecordQuery:%v", transactionID.ValidStart.UnixNano()), query._GetLogID())
}

func TestUnitTransactionRecordQueryNothingSet(t *testing.T) {
	query := NewTransactionRecordQuery()

	require.Equal(t, TransactionID{}, query.GetTransactionID())
	require.False(t, query.GetIncludeChildren())
	require.False(t, query.GetIncludeDuplicates())
	require.Empty(t, query.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 8*time.Second, query.GetMaxBackoff())
	require.Equal(t, 10, query.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, query.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, query.GetQueryPayment())
	require.Equal(t, Hbar{}, query.GetMaxQueryPayment())
}

func TestUnitTransactionRecordReceiptNotFound(t *testing.T) {
	responses := [][]interface{}{{
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status: services.ResponseCodeEnum_RECEIPT_NOT_FOUND,
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status: services.ResponseCodeEnum_RECEIPT_NOT_FOUND,
					},
				},
			},
		},
	}}
	client, server := NewMockClientAndServer(responses)
	defer server.Close()
	tx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(1)).
		Execute(client)
	client.SetMaxAttempts(2)
	require.NoError(t, err)
	record, err := tx.SetValidateStatus(true).GetRecord(client)
	require.Error(t, err)
	require.Equal(t, "exceptional precheck status RECEIPT_NOT_FOUND", err.Error())
	require.Equal(t, StatusReceiptNotFound, record.Receipt.Status)
}
