package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewAccountInfoQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 3}).
		SetQueryPaymentTransaction(mockTransaction)

	cupaloy.SnapshotT(t, query)
}

func TestAccountInfoQuery_Execute(t *testing.T) {
	operatorAccountID, err := AccountIDFromString(os.Getenv("OPERATOR_ID"))
	assert.NoError(t, err)

	operatorPrivateKey, err := Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	assert.NoError(t, err)

	client := ClientForTestnet().
		SetOperator(operatorAccountID, operatorPrivateKey)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

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

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, accountID, info.AccountID)
	assert.Equal(t, false, info.Deleted)
	assert.Equal(t, newKey.PublicKey(), info.Key)
	assert.Equal(t, newBalance.tinybar, info.Balance.tinybar)

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
