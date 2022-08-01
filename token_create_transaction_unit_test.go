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

func TestUnitTokenCreateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tokenCreate := NewTokenCreateTransaction().
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenCreate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenCreateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenCreate := NewTokenCreateTransaction().
		SetAutoRenewAccount(accountID).
		SetTreasuryAccountID(accountID)

	err = tokenCreate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenCreateTransactionGet(t *testing.T) {
	accountID := AccountID{Account: 7}

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetDecimals(3).
		SetCustomFees([]Fee{}).
		SetTokenType(TokenTypeFungibleCommon).
		SetSupplyType(TokenSupplyTypeFinite).
		SetMaxSupply(12312123123123).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(accountID).
		SetAdminKey(newKey).
		SetFreezeKey(newKey).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetPauseKey(newKey).
		SetExpirationTime(time.Now()).
		SetFreezeDefault(false).
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
	transaction.GetTokenMemo()
	transaction.GetDecimals()
	transaction.GetCustomFees()
	transaction.GetTokenType()
	transaction.GetSupplyType()
	transaction.GetMaxSupply()
	transaction.GetInitialSupply()
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

func TestUnitTokenCreateTransactionNothingSet(t *testing.T) {
	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewTokenCreateTransaction().
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
	transaction.GetTokenMemo()
	transaction.GetDecimals()
	transaction.GetCustomFees()
	transaction.GetTokenType()
	transaction.GetSupplyType()
	transaction.GetMaxSupply()
	transaction.GetInitialSupply()
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
	proto := transaction._Build().GetTokenCreation()
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
