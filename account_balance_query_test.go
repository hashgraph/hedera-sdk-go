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
	client := newTestClient(t, false)

	_, err := NewAccountBalanceQuery().
		SetAccountID(client.GetOperatorAccountID()).
		Execute(client)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_Execute(t *testing.T) {
	client := newTestClient(t, false)

	balance := NewAccountBalanceQuery().
		SetAccountID(client.GetOperatorAccountID())

	cost, err := balance.GetCost(client)
	assert.NoError(t, err)

	_, err = balance.SetMaxQueryPayment(cost).
		Execute(client)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_BigMax_Execute(t *testing.T) {
	client := newTestClient(t, false)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetAccountID(client.GetOperatorAccountID())

	cost, err := balance.GetCost(client)
	assert.NoError(t, err)

	_, err = balance.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_SmallMax_Execute(t *testing.T) {
	client := newTestClient(t, false)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetAccountID(client.GetOperatorAccountID())

	_, err := balance.GetCost(client)
	assert.NoError(t, err)

	_, err = balance.Execute(client)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_SetPayment_Execute(t *testing.T) {
	client := newTestClient(t, false)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetQueryPayment(NewHbar(0)).
		SetAccountID(client.GetOperatorAccountID())

	cost, err := balance.GetCost(client)
	assert.NoError(t, err)

	_, err = balance.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_SetPaymentOneTinybar_Execute(t *testing.T) {
	client := newTestClient(t, false)

	balance := NewAccountBalanceQuery().
		SetMaxQueryPayment(NewHbar(10000)).
		SetQueryPayment(NewHbar(0)).
		SetAccountID(client.GetOperatorAccountID())

	_, err := balance.GetCost(client)
	assert.NoError(t, err)

	_, err = balance.SetQueryPayment(HbarFromTinybar(1)).Execute(client)
	assert.NoError(t, err)
}

func Test_AccountBalance_NoAccount(t *testing.T) {
	client := newTestClient(t, false)

	_, err := NewAccountBalanceQuery().
		Execute(client)
	assert.Error(t, err)
	assert.True(t, err.Error() == "exceptional precheck status INVALID_ACCOUNT_ID")
}
