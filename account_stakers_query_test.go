package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountStakersQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	_, err := NewAccountStakersQuery().
		SetAccountID(client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.Error(t, err)
}
