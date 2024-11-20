//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenFeeScheduleUpdateTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1,
		DenominationTokenID: &tokenID,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.True(t, len(info.CustomFees) > 0)

}

func TestIntegrationTokenFeeScheduleUpdateTransactionWithFractional(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFractionalFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Numerator:        1,
		Denominator:      20,
		MinimumAmount:    1,
		MaximumAmount:    10,
		AssessmentMethod: FeeAssessmentMethodExclusive,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.True(t, len(info.CustomFees) > 0)

}

func TestIntegrationTokenFeeScheduleUpdateTransactionNoFeeScheduleKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetFeeScheduleKey(nil)
	})

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1,
		DenominationTokenID: &tokenID,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: TOKEN_HAS_NO_FEE_SCHEDULE_KEY", err.Error())
	}

}

func DisabledTestIntegrationTokenFeeScheduleUpdateTransactionWrongScheduleKey(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetFeeScheduleKey(newKey.PublicKey())
	})

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1,
		DenominationTokenID: &tokenID,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: INVALID_CUSTOM_FEE_SCHEDULE_KEY", err.Error())
	}

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.True(t, len(info.CustomFees) > 0)

}

func TestIntegrationTokenFeeScheduleUpdateTransactionScheduleAlreadyHasNoFees(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: CUSTOM_SCHEDULE_ALREADY_HAS_NO_FEES")

}

func TestIntegrationTokenFeeScheduleUpdateTransactionFractionalFeeOnlyForFungibleCommon(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createNft(&env)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFractionalFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Numerator:        1,
		Denominator:      20,
		MinimumAmount:    1,
		MaximumAmount:    10,
		AssessmentMethod: FeeAssessmentMethodExclusive,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_FRACTIONAL_FEE_ONLY_ALLOWED_FOR_FUNGIBLE_COMMON", err.Error())
	}

}

func TestIntegrationTokenFeeScheduleUpdateTransactionDenominationMustBeFungibleCommon(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	tokenIDNonFungible, err := createNft(&env)
	require.NoError(t, err)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              5,
		DenominationTokenID: &tokenIDNonFungible,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_FEE_DENOMINATION_MUST_BE_FUNGIBLE_COMMON", err.Error())
	}
}

func TestIntegrationTokenFeeScheduleUpdateTransactionCustomFeeListTooLong(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	resp, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1,
		DenominationTokenID: &tokenID,
	}

	feeArr := make([]Fee, 0)

	for i := 0; i < 21; i++ {
		feeArr = append(feeArr, customFee)
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees(feeArr).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_FEES_LIST_TOO_LONG", err.Error())
	}

}
