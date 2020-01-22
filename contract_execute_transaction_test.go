package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeContractExecuteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)


	parameters := NewContractFunctionParams().
		AddBytes([]byte{24, 43, 11})

	tx := NewContractExecuteTransaction().
		SetContractID(ContractID{Contract: 5}).
		SetGas(141).
		SetPayableAmount(HbarFromTinybar(10000)).
		SetFunction("someFunction", *parameters).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient).
		Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
