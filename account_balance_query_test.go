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
