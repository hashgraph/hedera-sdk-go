package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSerializeFileUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewFileUpdateTransaction().
		SetFileID(FileID{File: 5}).
		SetContents([]byte("there was a hole here")).
		SetExpirationTime(time.Unix(15415151511, 0)).
		AddKey(privateKey.PublicKey()).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)
	tx.Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
