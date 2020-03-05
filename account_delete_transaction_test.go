package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSerializeAccountDeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetDeleteAccountID(AccountID{Account: 3}).
		SetTransferAccountID(AccountID{Account: 2}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	cupaloy.SnapshotT(t, tx.String())
}

func TestAccountDeleteTransaction_Execute(t *testing.T) {
	operatorAccountID, err := AccountIDFromString(os.Getenv("OPERATOR_ID"))
	assert.NoError(t, err)

	operatorPrivateKey, err := Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	assert.NoError(t, err)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	client := ClientForTestnet().SetOperator(operatorAccountID, operatorPrivateKey)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	txID, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	accountID := receipt.GetAccountID()
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetDeleteAccountID(accountID).
		SetTransferAccountID(operatorAccountID).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(NewTransactionID(accountID)).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
