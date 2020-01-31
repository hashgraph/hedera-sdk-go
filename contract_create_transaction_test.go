package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSerializeContractCreateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewContractCreateTransaction().
		SetAdminKey(privateKey.PublicKey()).
		SetInitialBalance(HbarFromTinybar(1e3)).
		SetBytecodeFileID(FileID{File: 4}).
		SetGas(100).
		SetProxyAccountID(AccountID{Account: 3}).
		SetAutoRenewPeriod(60 * 60 * 24 * 14 * time.Second).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
