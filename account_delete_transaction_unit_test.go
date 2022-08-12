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
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountDeleteTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	deleteAccount := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(accountID)

	err = deleteAccount._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountDeleteTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	deleteAccount := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(accountID)

	err = deleteAccount._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountDeleteTransactionGet(t *testing.T) {
	spenderAccountID1 := AccountID{Account: 7}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTransferAccountID(spenderAccountID1).
		SetAccountID(spenderAccountID1).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetTransferAccountID()
	transaction.GetAccountID()
}

func TestUnitAccountDeleteTransactionSetNothing(t *testing.T) {
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	transaction.GetTransferAccountID()
	transaction.GetAccountID()
}

func TestUnitAccountDeleteTransactionProtoCheck(t *testing.T) {
	spenderAccountID1 := AccountID{Account: 7}
	transferAccountID := AccountID{Account: 8}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountDeleteTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTransferAccountID(transferAccountID).
		SetAccountID(spenderAccountID1).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction._Build().GetCryptoDelete()
	require.Equal(t, proto.TransferAccountID.String(), transferAccountID._ToProtobuf().String())
	require.Equal(t, proto.DeleteAccountID.String(), spenderAccountID1._ToProtobuf().String())
}
