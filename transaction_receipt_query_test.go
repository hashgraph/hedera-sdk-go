package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Receipt_Transaction(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(newBalance).
		FreezeWith(client)
	assert.NoError(t, err)

	tx, err = tx.SignWithOperator(client)
	assert.NoError(t, err)

	resp, err := tx.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	record, err := resp.GetRecord(client)
	assert.NoError(t, err)

	accountID := *record.Receipt.AccountID
	assert.NotNil(t, accountID)

	transaction, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
