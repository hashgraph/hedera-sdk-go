//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenCreateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionMultipleKeys(t *testing.T) {
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

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetFreezeKey(pubKeys[1]).
			SetWipeKey(pubKeys[2]).
			SetKycKey(pubKeys[3]).
			SetSupplyKey(pubKeys[4]).
			SetMetadataKey(pubKeys[5])
	})
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
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

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetFreezeKey(pubKeys[1]).
			SetWipeKey(pubKeys[2]).
			SetKycKey(pubKeys[3]).
			SetSupplyKey(pubKeys[4]).
			SetMetadataKey(pubKeys[5]).
			FreezeWith(env.Client)
		transaction.
			Sign(keys[0]).
			Sign(keys[1])
	})

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func DisabledTestIntegrationTokenNftCreateTransaction(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tokenID, err := createNft(&env)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionWithCustomFees(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
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
			})
	})

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionWithCustomFeesDenominatorZero(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
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
			})
	})
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: FRACTION_DIVIDES_BY_ZERO", err.Error())
	}
}

func TestIntegrationTokenCreateTransactionWithInvalidFeeCollectorAccountID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetCustomFees([]Fee{
				NewCustomFractionalFee().
					SetFeeCollectorAccountID(AccountID{}).
					SetNumerator(1).
					SetDenominator(20).
					SetMin(1).
					SetMax(10),
			})
	})
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_CUSTOM_FEE_COLLECTOR", err.Error())
	}
}

func TestIntegrationTokenCreateTransactionWithMaxLessThanMin(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
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
			})
	})
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: FRACTIONAL_FEE_MAX_AMOUNT_LESS_THAN_MIN_AMOUNT", err.Error())
	}
}

func TestIntegrationTokenCreateTransactionWithRoyaltyCustomFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
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
			})
	})
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenCreateTransactionWithRoyaltyCannotExceedOne(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetCustomFees([]Fee{
				NewCustomRoyaltyFee().
					SetFeeCollectorAccountID(env.OperatorID).
					SetNumerator(2).
					SetDenominator(1),
			})
	})
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: ROYALTY_FRACTION_CANNOT_EXCEED_ONE", err.Error())
	}
}

func TestIntegrationTokenCreateTransactionFeeCollectorMissing(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetCustomFees([]Fee{
				NewCustomRoyaltyFee().
					SetNumerator(1).
					SetDenominator(20),
			})
	})
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_CUSTOM_FEE_COLLECTOR", err.Error())
	}
}

func TestIntegrationTokenCreateTransactionRoyaltyFeeOnlyAllowedForNonFungibleUnique(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetCustomFees([]Fee{
				NewCustomRoyaltyFee().
					SetFeeCollectorAccountID(env.OperatorID).
					SetNumerator(1).
					SetDenominator(20),
			})
	})
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_ROYALTY_FEE_ONLY_ALLOWED_FOR_NON_FUNGIBLE_UNIQUE", err.Error())
	}
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

	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetTreasuryAccountID(accountID).
			SetAdminKey(newKey.PublicKey()).
			SetFreezeKey(newKey.PublicKey()).
			SetWipeKey(newKey.PublicKey()).
			SetKycKey(newKey.PublicKey()).
			SetSupplyKey(newKey.PublicKey()).
			FreezeWith(env.Client)

		transaction.Sign(newKey)
	})
	require.NoError(t, err)

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

func TestIntegrationTokenCreateTransactionMetadataKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	pubKey := newKey.PublicKey()

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(pubKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetMetadataKey(pubKey)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(2)).
		SetTokenID(tokenID).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)

	require.NoError(t, err)
	assert.Equal(t, pubKey, info.MetadataKey)

	err = CloseIntegrationTestEnv(env, &tokenID)
}
