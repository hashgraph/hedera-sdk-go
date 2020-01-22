package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContractCallQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	parameters := NewContractFunctionParams().
		AddBytes([]byte{24, 43, 11})

	query := NewContractCallQuery().
		SetGas(100).
		SetMaxResultSize(100).
		SetFunction("someFunction", *parameters).
		SetContractID(ContractID{Contract:3}).
		SetQueryPaymentTransaction(mockTransaction)

	cupaloy.SnapshotT(t, query.pb.String())
}

