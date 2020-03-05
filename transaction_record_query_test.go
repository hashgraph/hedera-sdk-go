package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTransactionRecordQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewTransactionRecordQuery().
		SetTransactionID(testTransactionID).
		SetQueryPaymentTransaction(mockTransaction)

	cupaloy.SnapshotT(t, query)
}
