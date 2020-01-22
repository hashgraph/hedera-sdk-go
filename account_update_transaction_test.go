package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeAccountUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewAccountUpdateTransaction().
		SetTransactionID(testTransactionID).
		SetAccountID(AccountID{Account: 3}).
		SetKey(privateKey.PublicKey()).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		Build(mockClient).
		Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
