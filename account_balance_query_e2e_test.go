//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountBalanceQueryCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountBalanceQuery().
		SetAccountID(env.OriginalOperatorID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountBalanceQueryCanGetTokenBalance(t *testing.T) {
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

	receipt, err := resp.GetReceipt(env.Client)
	require.NoError(t, err)

	tokenID := receipt.TokenID

	balance, err := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, balance.Tokens.Get(*tokenID), uint64(1000000))
	assert.Equal(t, balance.TokenDecimals.Get(*tokenID), uint64(3))

	err = CloseIntegrationTestEnv(env, tokenID)
	require.NoError(t, err)
}

func TestIntegrationAccountBalanceQueryGetCost(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.SetMaxQueryPayment(cost).
		Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountBalanceQuerySetBigMaxPayment(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountBalanceQuerySetSmallMaxPayment(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountBalanceQueryCanSetQueryPayment(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetQueryPayment(NewHbar(0)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.OperatorID)

	cost, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountBalanceQueryCostCanSetPaymentOneTinybar(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetQueryPayment(NewHbar(0)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := balance.GetCost(env.Client)
	require.NoError(t, err)

	_, err = balance.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func TestIntegrationAccountBalanceQueryNoAccountIDError(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	assert.True(t, err.Error() == "exceptional precheck status INVALID_ACCOUNT_ID")

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
