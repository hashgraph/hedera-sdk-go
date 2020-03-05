package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGetBySolidityIDQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewGetBySolidityIDQuery().
		SetSolidityID("not a real solidity id").
		SetQueryPaymentTransaction(mockTransaction)

	cupaloy.SnapshotT(t, query)
}
