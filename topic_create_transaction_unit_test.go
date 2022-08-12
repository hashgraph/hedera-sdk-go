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

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTopicCreateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	topicCreate := NewTopicCreateTransaction().
		SetAutoRenewAccountID(accountID)

	err = topicCreate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTopicCreateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	topicCreate := NewTopicCreateTransaction().
		SetAutoRenewAccountID(accountID)

	err = topicCreate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTopicCreateTransactionGet(t *testing.T) {
	accountID := AccountID{Account: 3}

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTopicCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAutoRenewAccountID(accountID).
		SetAdminKey(newKey).
		SetSubmitKey(newKey).
		SetTopicMemo("ad").
		SetAutoRenewPeriod(60 * time.Second).
		SetMaxTransactionFee(NewHbar(10)).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetRegenerateTransactionID(false).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetAutoRenewAccountID()
	transaction.GetAdminKey()
	transaction.GetSubmitKey()
	transaction.GetTopicMemo()
	transaction.GetAutoRenewPeriod()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTopicCreateTransactionNothingSet(t *testing.T) {
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTopicCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetAutoRenewAccountID()
	transaction.GetAdminKey()
	transaction.GetSubmitKey()
	transaction.GetTopicMemo()
	transaction.GetAutoRenewPeriod()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTopicCreateTransactionProtoCheck(t *testing.T) {
	accountID := AccountID{Account: 23}
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	newKey2, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transaction, err := NewTopicCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(newKey).
		SetSubmitKey(newKey2).
		SetAutoRenewAccountID(accountID).
		SetTopicMemo("memo").
		SetAutoRenewPeriod(time.Second * 3).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction._Build().GetConsensusCreateTopic()
	require.Equal(t, proto.AdminKey.String(), newKey._ToProtoKey().String())
	require.Equal(t, proto.SubmitKey.String(), newKey2._ToProtoKey().String())
	require.Equal(t, proto.Memo, "memo")
	require.Equal(t, proto.AutoRenewPeriod.Seconds, _DurationToProtobuf(time.Second*3).Seconds)
	require.Equal(t, proto.AutoRenewAccount.String(), accountID._ToProtobuf().String())
}
