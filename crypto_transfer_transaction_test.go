package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSerializeCryptoTransferTransaction(t *testing.T) {
	tx, err := newMockTransaction()
	assert.NoError(t, err)

	cupaloy.SnapshotT(t, tx.String())
}

func TestCryptoTransferTransaction_Execute(t *testing.T) {
	operatorAccountID, err := AccountIDFromString(os.Getenv("OPERATOR_ID"))
	assert.NoError(t, err)

	operatorPrivateKey, err := Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	assert.NoError(t, err)

	client := ClientForTestnet().
		SetOperator(operatorAccountID, operatorPrivateKey)

	txID, err := NewCryptoTransferTransaction().
		AddSender(operatorAccountID, NewHbar(1)).
		AddRecipient(AccountID{Account: 3}, NewHbar(1)).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
