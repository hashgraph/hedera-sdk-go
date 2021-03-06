package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountBalanceQuery_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance, err := NewAccountBalanceQuery().
		SetAccountID(env.OriginalOperatorID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.NoError(t, err)

	println(balance.Hbars.String())

	_, err = NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestAccountBalanceQuery_TokenBalance(t *testing.T) {
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
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := receipt.TokenID

	balance, err := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, balance.Tokens.Get(*tokenID), uint64(1000000))
	assert.Equal(t, balance.TokenDecimals.Get(*tokenID), uint64(3))

	err = CloseIntegrationTestEnv(env, tokenID)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := balance.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = balance.SetMaxQueryPayment(cost).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_BigMax_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := balance.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = balance.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_SmallMax_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := balance.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = balance.Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_SetPayment_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetQueryPayment(NewHbar(0)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := balance.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = balance.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_SetPaymentOneTinybar_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetQueryPayment(NewHbar(0)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := balance.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = balance.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func Test_AccountBalance_NoAccount(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	assert.True(t, err.Error() == "exceptional precheck status INVALID_ACCOUNT_ID")

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
