package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContractInfoQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewContractInfoQuery().
		SetContractID(ContractID{Contract: 3}).
		SetQueryPaymentTransaction(mockTransaction)

	cupaloy.SnapshotT(t, query.pb.String())
}
