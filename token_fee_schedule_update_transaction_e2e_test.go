//go:build all || e2e
// +build all e2e

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
    "fmt"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenFeeScheduleUpdateTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
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

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.True(t, len(info.CustomFees) > 0)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTransactionWithFractional(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
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

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)
	assert.True(t, len(info.CustomFees) > 0)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTransactionNoFeeScheduleKey(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
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

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
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

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: TOKEN_HAS_NO_FEE_SCHEDULE_KEY", err.Error())
	}

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func DisabledTestIntegrationTokenFeeScheduleUpdateTransactionWrongScheduleKey(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(newKey.PublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
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

	_, err = resp.GetReceipt(env.Client)
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

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTransactionScheduleAlreadyHasNoFees(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTransactionFractionalFeeOnlyForFungibleCommon(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetSupplyType(TokenSupplyTypeFinite).
		SetMaxSupply(5).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
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

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_FRACTIONAL_FEE_ONLY_ALLOWED_FOR_FUNGIBLE_COMMON", err.Error())
	}

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTransactionDenominationMustBeFungibleCommon(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(6).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetSupplyType(TokenSupplyTypeFinite).
		SetMaxSupply(6).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenIDNonFungible := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
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

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_FEE_DENOMINATION_MUST_BE_FUNGIBLE_COMMON", err.Error())
	}
}

func TestIntegrationTokenFeeScheduleUpdateTransactionCustomFeeListTooLong(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
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

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: CUSTOM_FEES_LIST_TOO_LONG", err.Error())
	}

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTransactionCanExecute2(t *testing.T) {
	env := NewIntegrationTestEnv(t)

    fmt.Printf("Client operator: %v\n", env.Client.GetOperatorAccountID().String())

    key1, err := PrivateKeyGenerateEd25519()
    require.NoError(t, err)
    key2, err := PrivateKeyGenerateEd25519()
    require.NoError(t, err)
    key3, err := PrivateKeyGenerateEd25519()
    require.NoError(t, err)
    key4, err := PrivateKeyGenerateEd25519()
    require.NoError(t, err)

    resp, err := NewAccountCreateTransaction().
        SetKey(key1).
        SetInitialBalance(NewHbar(1)).
        Execute(env.Client)
    require.NoError(t, err)

    receipt, err := resp.GetReceipt(env.Client)
    require.NoError(t, err)

    accountID1 := *receipt.AccountID
    fmt.Printf("AccountID 1: %v\n", accountID1.String())

    resp, err = NewAccountCreateTransaction().
        SetKey(key2).
        SetInitialBalance(NewHbar(1)).
        Execute(env.Client)
    require.NoError(t, err)

    receipt, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    accountID2 := *receipt.AccountID
    fmt.Printf("AccountID 2: %v\n", accountID2.String())

    resp, err = NewAccountCreateTransaction().
        SetKey(key3).
        SetInitialBalance(NewHbar(1)).
        Execute(env.Client)
    require.NoError(t, err)

    receipt, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    accountID3 := *receipt.AccountID
    fmt.Printf("AccountID 3 (Treasury): %v\n", accountID3.String())

    resp, err = NewAccountCreateTransaction().
        SetKey(key4).
        SetInitialBalance(NewHbar(1)).
        Execute(env.Client)
    require.NoError(t, err)

    receipt, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    accountID4 := *receipt.AccountID
    fmt.Printf("AccountID 4: %v\n", accountID4.String())

	customFee1 := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &accountID1,
		},
		Amount:              1,
		DenominationTokenID: nil,
	}

    fmt.Printf("Created fixed custom fees of 1 to collector %v\n", accountID1.String())

	customFee2 := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &accountID2,
		},
		Amount:              1,
		DenominationTokenID: nil,
	}

    fmt.Printf("Created fixed custom fees of 1 to collector %v\n", accountID2.String())

	customFee4 := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &accountID4,
		},
		Amount:              1,
		DenominationTokenID: nil,
	}

    fmt.Printf("Created fixed custom fees of 1 to collector %v\n", accountID4.String())

    customFees := []Fee{customFee1, customFee2, customFee4}

	tokenCreate, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
        SetTreasuryAccountID(accountID3).
        SetAdminKey(env.Client.GetOperatorPublicKey()).
        SetCustomFees(customFees).
        FreezeWith(env.Client)
	require.NoError(t, err)

    tokenCreate.Sign(key3)

    resp, err = tokenCreate.Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

    associate, err := NewTokenAssociateTransaction().
        SetAccountID(accountID1).
        SetTokenIDs(tokenID).
        FreezeWith(env.Client)
    require.NoError(t, err)

    associate.Sign(key1)

    resp, err = associate.Execute(env.Client)

    _, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    fmt.Printf("Associated account ID 1 %v with token\n", accountID1.String())

    associate, err = NewTokenAssociateTransaction().
        SetAccountID(accountID2).
        SetTokenIDs(tokenID).
        FreezeWith(env.Client)
    require.NoError(t, err)

    associate.Sign(key2)

    resp, err = associate.Execute(env.Client)

    _, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    fmt.Printf("Associated account ID 2 %v with token\n", accountID2.String())

    associate, err = NewTokenAssociateTransaction().
        SetAccountID(accountID4).
        SetTokenIDs(tokenID).
        FreezeWith(env.Client)
    require.NoError(t, err)

    associate.Sign(key4)

    resp, err = associate.Execute(env.Client)

    _, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    fmt.Printf("Associated account ID 4 %v with token\n", accountID4.String())

    resp, err = NewTokenAssociateTransaction().
        SetAccountID(env.Client.GetOperatorAccountID()).
        SetTokenIDs(tokenID).
        Execute(env.Client)
    require.NoError(t, err)

    _, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    fmt.Printf("Associated operator account %v with token\n", env.Client.GetOperatorAccountID().String())

    transfer, err := NewTransferTransaction().
        AddTokenTransfer(tokenID, accountID3, -10).
        AddTokenTransfer(tokenID, accountID1, 10).
        FreezeWith(env.Client)
    require.NoError(t, err)

    transfer.Sign(key3)

    resp, err = transfer.Execute(env.Client)
    require.NoError(t, err)

    _, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    fmt.Printf("Transfered 10 token from treasury to account ID %v\n", accountID1.String())

    transfer, err = NewTransferTransaction().
        AddTokenTransfer(tokenID, accountID1, -8).
        AddTokenTransfer(tokenID, accountID2, 8).
        FreezeWith(env.Client)
    require.NoError(t, err)

    transfer.Sign(key1)

    resp, err = transfer.Execute(env.Client)
    require.NoError(t, err)

    _, err = resp.GetReceipt(env.Client)
    require.NoError(t, err)

    fmt.Printf("Transfered 8 token from %v to %v\n", accountID1.String(), accountID2.String())

    record, err := resp.GetRecord(env.Client)
    require.NoError(t, err)

    for tokenID, transfers := range record.TokenTransfers {
        fmt.Printf("%v: %+v\n", tokenID.String(), transfers)
    }

    for _, fee := range record.AssessedCustomFees {
        fmt.Printf("%v\n", fee.String())
    }
}
