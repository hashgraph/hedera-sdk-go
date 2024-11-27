//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTokenInfoQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetTokenMetadata([]byte{1, 2, 3}).
			SetKycKey(env.Client.GetOperatorPublicKey()).
			SetDecimals(3)
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
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
	assert.Equal(t, info.MetadataKey.String(), env.Client.GetOperatorPublicKey().String())
	assert.Equal(t, info.Metadata, []byte{1, 2, 3})
	assert.False(t, *info.DefaultFreezeStatus)
	assert.False(t, *info.DefaultKycStatus)

}

func TestIntegrationTokenInfoQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	infoQuery := NewTokenInfoQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(env.Client)
	require.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTokenInfoQuerySetBigMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	infoQuery := NewTokenInfoQuery().
		SetMaxQueryPayment(NewHbar(1000000)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(env.Client)
	require.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTokenInfoQuerySetSmallMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	infoQuery := NewTokenInfoQuery().
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetTokenID(tokenID)

	cost, err := infoQuery.GetCost(env.Client)
	require.NoError(t, err)

	_, err = infoQuery.Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "cost of TokenInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tℏ", err.Error())
	}

}

func TestIntegrationTokenInfoQueryInsufficientCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	infoQuery := NewTokenInfoQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetTokenID(tokenID)

	_, err = infoQuery.GetCost(env.Client)
	require.NoError(t, err)

	_, err = infoQuery.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INSUFFICIENT_TX_FEE", err.Error())
	}

}

func TestIntegrationTokenInfoQueryNoPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	tokenID, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetDecimals(3).
			SetKycKey(env.Client.GetOperatorPublicKey())
	})
	require.NoError(t, err)

	info, err := NewTokenInfoQuery().
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

}

func TestIntegrationTokenInfoQueryNoTokenID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewTokenInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INVALID_TOKEN_ID", err.Error())
	}

}
