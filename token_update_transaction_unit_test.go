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

func TestUnitTokenUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tokenUpdate := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenUpdate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenUpdate := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenUpdate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenUpdateTransactionGet(t *testing.T) {
	accountID := AccountID{Account: 7}

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetTreasuryAccountID(accountID).
		SetTokenID(TokenID{Token: 7}).
		SetAdminKey(newKey).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetPauseKey(newKey).
		SetExpirationTime(time.Now()).
		SetAutoRenewPeriod(60 * time.Second).
		SetAutoRenewAccount(accountID).
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

	transaction.GetTokenName()
	transaction.GetTokenSymbol()
	transaction.GetTreasuryAccountID()
	transaction.GetAdminKey()
	transaction.GetFreezeKey()
	transaction.GetWipeKey()
	transaction.GetKycKey()
	transaction.GetSupplyKey()
	transaction.GetPauseKey()
	transaction.GetExpirationTime()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
}

func TestUnitTokenUpdateTransactionNothingSet(t *testing.T) {
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetTokenName()
	transaction.GetTokenSymbol()
	transaction.GetTreasuryAccountID()
	transaction.GetAdminKey()
	transaction.GetFreezeKey()
	transaction.GetWipeKey()
	transaction.GetKycKey()
	transaction.GetSupplyKey()
	transaction.GetPauseKey()
	transaction.GetExpirationTime()
	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetMaxTransactionFee()
	transaction.GetRegenerateTransactionID()
	proto := transaction._Build().GetTokenUpdate()
	require.Nil(t, proto.Token)
	require.Nil(t, proto.AutoRenewPeriod)
	require.Nil(t, proto.AutoRenewAccount)
	require.Nil(t, proto.AdminKey)
	require.Nil(t, proto.Expiry)
	require.Nil(t, proto.FeeScheduleKey)
	require.Nil(t, proto.FreezeKey)
	require.Nil(t, proto.KycKey)
	require.Nil(t, proto.FeeScheduleKey)
	require.Nil(t, proto.PauseKey)
	require.Nil(t, proto.SupplyKey)
	require.Nil(t, proto.Treasury)
}

func TestUnitTokenUpdateTransactionKeyCheck(t *testing.T) {
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	keys := make([]PrivateKey, 7)
	var err error

	for i := 0; i < 7; i++ {
		keys[i], err = PrivateKeyGenerateEd25519()
		require.NoError(t, err)
	}

	transaction, err := NewTokenUpdateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetAdminKey(keys[0]).
		SetFreezeKey(keys[1]).
		SetWipeKey(keys[2]).
		SetKycKey(keys[3]).
		SetSupplyKey(keys[4]).
		SetPauseKey(keys[5]).
		SetFeeScheduleKey(keys[6]).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction._Build().GetTokenUpdate()
	require.Equal(t, proto.AdminKey.String(), keys[0]._ToProtoKey().String())
	require.Equal(t, proto.FreezeKey.String(), keys[1]._ToProtoKey().String())
	require.Equal(t, proto.WipeKey.String(), keys[2]._ToProtoKey().String())
	require.Equal(t, proto.KycKey.String(), keys[3]._ToProtoKey().String())
	require.Equal(t, proto.SupplyKey.String(), keys[4]._ToProtoKey().String())
	require.Equal(t, proto.PauseKey.String(), keys[5]._ToProtoKey().String())
	require.Equal(t, proto.FeeScheduleKey.String(), keys[6]._ToProtoKey().String())
}
