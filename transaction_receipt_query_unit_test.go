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

func TestUnitTransactionReceiptQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	transactionID := TransactionIDGenerate(accountID)
	require.NoError(t, err)

	receiptQuery := NewTransactionReceiptQuery().
		SetTransactionID(transactionID)

	err = receiptQuery._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransactionReceiptQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	transactionID := TransactionIDGenerate(accountID)
	require.NoError(t, err)

	receiptQuery := NewTransactionReceiptQuery().
		SetTransactionID(transactionID)

	err = receiptQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTransactionReceiptQueryGet(t *testing.T) {
	txID := TransactionIDGenerate(AccountID{Account: 7})

	balance := NewTransactionReceiptQuery().
		SetTransactionID(txID).
		SetIncludeDuplicates(true).
		SetIncludeChildren(true).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetTransactionID()
	balance.GetIncludeChildren()
	balance.GetIncludeDuplicates()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitTransactionReceiptQueryNothingSet(t *testing.T) {
	balance := NewTransactionReceiptQuery()

	balance.GetTransactionID()
	balance.GetIncludeChildren()
	balance.GetIncludeDuplicates()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}
func TestUnitTransactionReceiptNotFound(t *testing.T) {
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
	receipt, err := tx.SetValidateStatus(true).GetReceipt(client)
	require.Error(t, err)
	require.Equal(t, "exceptional precheck status RECEIPT_NOT_FOUND", err.Error())
	require.Equal(t, StatusReceiptNotFound, receipt.Status)
}
