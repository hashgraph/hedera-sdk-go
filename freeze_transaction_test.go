package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeFreezeTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewFreezeTransaction().
		SetTransactionID(testTransactionID).
		SetStartTime(12, 30).
		SetEndTime(14, 30).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
