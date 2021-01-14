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
	client := newTestClient(t)

	_, err := NewAccountBalanceQuery().
		SetAccountID(client.GetOperatorAccountID()).
		Execute(client)
	assert.NoError(t, err)
}

func TestAccountBalanceQueryCost_Execute(t *testing.T) {
	client := newTestClient(t)

	balance := NewAccountBalanceQuery().
		SetAccountID(client.GetOperatorAccountID())

	cost, err := balance.GetCost(client)
	assert.NoError(t, err)

	_, err = balance.SetMaxQueryPayment(cost).
		Execute(client)
	assert.NoError(t, err)
}

func Test_AccountBalance_NoAccount(t *testing.T) {
	client := newTestClient(t)

	_, err := NewAccountBalanceQuery().
		Execute(client)
	assert.Error(t, err)
	assert.True(t, err.Error() == "exceptional precheck status INVALID_ACCOUNT_ID")
}
