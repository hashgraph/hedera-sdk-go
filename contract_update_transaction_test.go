package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSerializeContractUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewContractUpdateTransaction().
		SetContractID(ContractID{Contract: 3}).
		SetAdminKey(privateKey.PublicKey()).
		SetBytecodeFileID(FileID{File: 5}).
		SetExpirationTime(time.Unix(1569375111277, 0)).
		SetProxyAccountID(AccountID{Account: 3}).
		SetAutoRenewPeriod(60 * 60 * 24 * 14 * time.Second).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient).
		Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}
