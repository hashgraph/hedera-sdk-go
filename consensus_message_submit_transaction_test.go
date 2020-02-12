package hedera

import (
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestSerializeConsensusMessageSubmitTransaction(t *testing.T) {
	date := time.Unix(1554158542, 0)

	testTopicID := ConsensusTopicID{Topic: 99}

	key, err := Ed25519PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	tx, err := NewConsensusMessageSubmitTransaction().
		SetTopicID(testTopicID).
		SetMessage([]byte("Hello Hashgraph")).
		SetTransactionValidDuration(24 * time.Hour).
		SetNodeAccountID(AccountID{Account: 3}).
		SetTransactionID(TransactionID{
			AccountID:  AccountID{Account: 2},
			ValidStart: date,
		}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		Build(nil)

	assert.NoError(t, err)

	tx.Sign(key)

	cupaloy.SnapshotT(t, tx.String())
}
