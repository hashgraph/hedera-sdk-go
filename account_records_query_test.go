package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccountRecordQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewAccountRecordsQuery().
		SetAccountID(AccountID{Account: 3}).
		SetQueryPaymentTransaction(mockTransaction)

	cupaloy.SnapshotT(t, query)
}
