package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccountBalanceQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewAccountBalanceQuery().
		SetAccountID(AccountID{Account:3}).
		SetQueryPaymentTransaction(mockTransaction)

	cupaloy.SnapshotT(t, query.pb.String())
}
