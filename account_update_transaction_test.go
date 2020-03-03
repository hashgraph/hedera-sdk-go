package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSerializeAccountUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetTransactionID(testTransactionID).
		SetAccountID(AccountID{Account: 3}).
		SetKey(privateKey.PublicKey()).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	cupaloy.SnapshotT(t, tx)
}

func TestAccountUpdateTransaction_Execute(t *testing.T) {
	operatorAccountID, err := AccountIDFromString(os.Getenv("OPERATOR_ID"))
	assert.NoError(t, err)

	operatorPrivateKey, err := Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	assert.NoError(t, err)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newKey2, err := GenerateEd25519PrivateKey()
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

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetKey(newKey2.PublicKey()).
		SetMaxTransactionFee(NewHbar(1)).
		Build(client)
	assert.NoError(t, err)

	tx.Sign(newKey)
	tx.Sign(newKey2)

	_, err = tx.Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey(), info.Key)

	tx, err = NewAccountDeleteTransaction().
		SetDeleteAccountID(accountID).
		SetTransferAccountID(operatorAccountID).
		SetMaxTransactionFee(NewHbar(1)).
		Build(client)
	assert.NoError(t, err)

	tx.Sign(newKey2)

	txID, err = tx.Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)

	assert.NoError(t, err)
}
