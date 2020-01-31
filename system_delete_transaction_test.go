package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSerializeSystemDeleteFileIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewSystemDeleteTransaction().
		SetFileID(FileID{File: 3}).
		SetExpirationTime(time.Unix(15415151511, 0)).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}

func TestSerializeSystemDeleteContractIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewSystemDeleteTransaction().
	 	SetContractID(ContractID{Contract: 3}).
		SetExpirationTime(time.Unix(15415151511, 0)).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)
	tx.Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
