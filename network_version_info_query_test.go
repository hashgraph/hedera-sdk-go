package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkVersionInfoQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	_, err := NewNetworkVersionQuery().Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED"), err.Error())
	}
}
