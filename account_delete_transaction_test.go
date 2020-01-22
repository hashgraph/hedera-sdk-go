package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeAccountDeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewAccountDeleteTransaction().
		SetDeleteAccountID(AccountID{Account: 3}).
		SetTransferAccountID(AccountID{Account: 2}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient).
		Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
