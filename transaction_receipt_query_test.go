package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReceiptQueryTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		FreezeWith(client)
	assert.NoError(t, err)

	tx.SignWithOperator(client)

	resp, err := tx.Execute(client)

	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	record, err := resp.GetRecord(client)
	assert.NoError(t, err)

	accountID := *record.Receipt.AccountID
	assert.NotNil(t, accountID)

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	transcation, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs(nodeIDs).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transcation.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
