package hedera

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeAccountBalanceQuery(t *testing.T) {
	query := NewAccountBalanceQuery().
		SetAccountID(AccountID{Account: 3}).
		Query

	assert.Equal(t, `cryptogetAccountBalance:{header:{}accountID:{accountNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestAccountBalanceQuery_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
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
}

func Test_AccountBalance_NoAccount(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountBalanceQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	assert.True(t, err.Error() == "exceptional precheck status INVALID_ACCOUNT_ID")
}
