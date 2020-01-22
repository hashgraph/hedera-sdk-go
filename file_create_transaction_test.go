package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeFileCreateTransaction(t *testing.T) {
	date := time.Unix(1554158542, 0)

	key, err := Ed25519PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	tx := NewFileCreateTransaction().
		AddKey(key.PublicKey()).
		SetContents([]byte{1, 2, 3, 4}).
		SetExpirationTime(date).
		SetNodeAccountID(AccountID{Account: 3}).
		SetTransactionID(TransactionID{
			AccountID:  AccountID{Account: 2},
			ValidStart: date,
		}).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		Build(nil).
		Sign(key)

	cupaloy.SnapshotT(t, tx.String())
}
