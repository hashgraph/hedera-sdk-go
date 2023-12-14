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

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenInfoQueryCanExecute(t *testing.T) {
	t.Parallel()
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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	info, err := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(2)).
		SetTokenID(tokenID).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, info.TokenID, tokenID)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(3))
	assert.Equal(t, info.Treasury, env.Client.GetOperatorAccountID())
	assert.NotNil(t, info.AdminKey)
	assert.NotNil(t, info.KycKey)
	assert.NotNil(t, info.FreezeKey)
	assert.NotNil(t, info.WipeKey)
	assert.NotNil(t, info.SupplyKey)
	assert.Equal(t, info.AdminKey.String(), env.Client.GetOperatorPublicKey().String())
	assert.Equal(t, info.KycKey.String(), env.Client.GetOperatorPublicKey().String())
	assert.Equal(t, info.FreezeKey.String(), env.Client.GetOperatorPublicKey().String())
	assert.Equal(t, info.WipeKey.String(), env.Client.GetOperatorPublicKey().String())
	assert.Equal(t, info.SupplyKey.String(), env.Client.GetOperatorPublicKey().String())
	assert.False(t, *info.DefaultFreezeStatus)
	assert.False(t, *info.DefaultKycStatus)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenInfoQueryGetCost(t *testing.T) {
	t.Parallel()
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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	infoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(env.Client)
	require.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenInfoQuerySetBigMaxPayment(t *testing.T) {
	t.Parallel()
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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	infoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1000000)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(env.Client)
	require.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenInfoQuerySetSmallMaxPayment(t *testing.T) {
	t.Parallel()
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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	infoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(env.Client)
	require.NoError(t, err)

	_, err = infoQuery.Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "cost of TokenInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tℏ", err.Error())
	}

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenInfoQueryInsufficientCost(t *testing.T) {
	t.Parallel()
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

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	infoQuery := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		SetTokenID(tokenID)

	_, err = infoQuery.GetCost(env.Client)
	require.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INSUFFICIENT_TX_FEE", err.Error())
	}

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenInfoQueryNoPayment(t *testing.T) {
	t.Parallel()
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
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := *receipt.TokenID

	info, err := NewTokenInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		SetTokenID(tokenID).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, info.TokenID, tokenID)
	assert.Equal(t, info.Name, "ffff")
	assert.Equal(t, info.Symbol, "F")
	assert.Equal(t, info.Decimals, uint32(3))
	assert.Equal(t, info.Treasury, env.Client.GetOperatorAccountID())
	assert.False(t, *info.DefaultFreezeStatus)
	assert.False(t, *info.DefaultKycStatus)

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationTokenInfoQueryNoTokenID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	_, err := NewTokenInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INVALID_TOKEN_ID", err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
