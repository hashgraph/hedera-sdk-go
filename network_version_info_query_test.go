package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkVersionInfoQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	_, err := NewNetworkVersionQuery().
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED"), err.Error())
	}
}

func TestNetworkVersionInfoQueryCost_Execute(t *testing.T) {
	client := newTestClient(t)

	query := NewNetworkVersionQuery()

	cost, err := query.GetCost(client)
	assert.Error(t, err)

	_, err = query.SetQueryPayment(cost).Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED"), err.Error())
	}
}
