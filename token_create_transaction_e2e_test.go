//go:build all || e2e
// +build all e2e

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

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenCreateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionMultipleKeys(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 5)
	pubKeys := make([]PublicKey, 5)

	for i := range keys {
		newKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(pubKeys[0]).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(pubKeys[1]).
		SetWipeKey(pubKeys[2]).
		SetKycKey(pubKeys[3]).
		SetSupplyKey(pubKeys[4]).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionNoKeys(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 6)
	pubKeys := make([]PublicKey, 6)

	for i := range keys {
		newKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(pubKeys[0]).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	info, err := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(env.Client)

	require.NoError(t, err)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(0))
	assert.Equal(t, info.TotalSupply, uint64(0))
	assert.Equal(t, info.Treasury.String(), env.Client.GetOperatorAccountID().String())
	assert.Nil(t, info.AdminKey)
	assert.Nil(t, info.FreezeKey)
	assert.Nil(t, info.KycKey)
	assert.Nil(t, info.WipeKey)
	assert.Nil(t, info.SupplyKey)
	assert.Nil(t, info.DefaultFreezeStatus)
	assert.Nil(t, info.DefaultKycStatus)
	assert.NotNil(t, info.AutoRenewPeriod)
	assert.Equal(t, *info.AutoRenewPeriod, 7890000*time.Second)
	assert.NotNil(t, info.AutoRenewAccountID)
	assert.Equal(t, info.AutoRenewAccountID.String(), env.Client.GetOperatorAccountID().String())
	assert.NotNil(t, info.ExpirationTime)
}

func TestIntegrationTokenCreateTransactionAdminSign(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 6)
	pubKeys := make([]PublicKey, 6)

	for i := range keys {
		newKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(pubKeys[0]).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenCreate, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(pubKeys[1]).
		SetWipeKey(pubKeys[2]).
		SetKycKey(pubKeys[3]).
		SetSupplyKey(pubKeys[4]).
		SetFreezeDefault(false).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tokenCreate.
		Sign(keys[0]).
		Sign(keys[1]).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	assert.NotNil(t, receipt.TokenID)

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func DisabledTestIntegrationTokenNftCreateTransaction(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetSupplyType(TokenSupplyTypeFinite).
		SetMaxSupply(5).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionWithCustomFees(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]Fee{
			NewCustomFixedFee().
				SetFeeCollectorAccountID(env.OperatorID).
				SetAmount(10),
			NewCustomFractionalFee().
				SetFeeCollectorAccountID(env.OperatorID).
				SetNumerator(1).
				SetDenominator(20).
				SetMin(1).
				SetAssessmentMethod(true).
				SetMax(10),
		}).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionWithCustomFeesDenominatorZero(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]Fee{
			CustomFixedFee{
				CustomFee: CustomFee{
					FeeCollectorAccountID: &env.OperatorID,
				},
				Amount: 10,
			},
			CustomFractionalFee{
				CustomFee: CustomFee{
					FeeCollectorAccountID: &env.OperatorID,
				},
				Numerator:     1,
				Denominator:   0,
				MinimumAmount: 1,
				MaximumAmount: 10,
			},
		}).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: FRACTION_DIVIDES_BY_ZERO", err.Error())
	}
}

func TestIntegrationTokenCreateTransactionWithInvalidFeeCollectorAccountID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]Fee{
			NewCustomFractionalFee().
				SetFeeCollectorAccountID(AccountID{}).
				SetNumerator(1).
				SetDenominator(20).
				SetMin(1).
				SetMax(10),
		}).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_CUSTOM_FEE_COLLECTOR", err.Error())
	}

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionWithMaxLessThanMin(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]Fee{
			CustomFractionalFee{
				CustomFee: CustomFee{
					FeeCollectorAccountID: &env.OperatorID,
				},
				Numerator:     1,
				Denominator:   20,
				MinimumAmount: 100,
				MaximumAmount: 10,
			},
		}).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: FRACTIONAL_FEE_MAX_AMOUNT_LESS_THAN_MIN_AMOUNT", err.Error())
	}

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionWithRoyaltyCustomFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]Fee{
			NewCustomRoyaltyFee().
				SetFeeCollectorAccountID(env.OperatorID).
				SetNumerator(1).
				SetDenominator(20).
				SetFallbackFee(
					NewCustomFixedFee().
						SetFeeCollectorAccountID(env.OperatorID).
						SetAmount(10),
				),
		}).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionWithRoyaltyCannotExceedOne(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]Fee{
			NewCustomRoyaltyFee().
				SetFeeCollectorAccountID(env.OperatorID).
				SetNumerator(2).
				SetDenominator(1),
		}).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: ROYALTY_FRACTION_CANNOT_EXCEED_ONE", err.Error())
	}

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionFeeCollectorMissing(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]Fee{
			NewCustomRoyaltyFee().
				SetNumerator(1).
				SetDenominator(20),
		}).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_CUSTOM_FEE_COLLECTOR", err.Error())
	}

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionRoyaltyFeeOnlyAllowedForNonFungibleUnique(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetTokenType(TokenTypeFungibleCommon).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetCustomFees([]Fee{
			NewCustomRoyaltyFee().
				SetFeeCollectorAccountID(env.OperatorID).
				SetNumerator(1).
				SetDenominator(20),
		}).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_ROYALTY_FEE_ONLY_ALLOWED_FOR_NON_FUNGIBLE_UNIQUE", err.Error())
	}

	err = CloseIntegrationTestEnv(env, receipt.TokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenAccountStillOwnsNfts(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tokTx, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenMemo("fnord").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetSupplyType(TokenSupplyTypeFinite).
		SetMaxSupply(5).
		SetTreasuryAccountID(accountID).
		SetAdminKey(newKey.PublicKey()).
		SetFreezeKey(newKey.PublicKey()).
		SetWipeKey(newKey.PublicKey()).
		SetKycKey(newKey.PublicKey()).
		SetSupplyKey(newKey.PublicKey()).
		SetFreezeDefault(false).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokTx.Sign(newKey)

	resp, err = tokTx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID
	metaData := make([]byte, 50, 101)

	mintTx, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadata(metaData).
		FreezeWith(env.Client)
	require.NoError(t, err)

	mintTx.Sign(newKey)

	mint, err := mintTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	deleteTx, err := NewTokenDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenID(tokenID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	deleteTx.Sign(newKey)

	resp, err = deleteTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}
