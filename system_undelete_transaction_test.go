package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeSystemUndeleteFileIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewSystemUndeleteTransaction().
		SetFileID(FileID{File: 3}).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		SetTransactionID(testTransactionID).
		Build(mockClient).
		Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}

func TestSerializeSystemUndeleteContractIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewSystemUndeleteTransaction().
		SetContractID(ContractID{Contract: 3}).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		SetTransactionID(testTransactionID).
		Build(mockClient).
		Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
