package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Receipt_Transaction(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	tx, err = tx.SignWithOperator(env.Client)
	assert.NoError(t, err)

	resp, err := tx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetRecord(env.Client)
	assert.NoError(t, err)

	// accountID := *record.Receipt.AccountID
	// assert.NotNil(t, accountID)

	// transaction, err := NewAccountDeleteTransaction().
	// 	SetNodeAccountIDs([]AccountID{resp.NodeID}).
	// 	SetAccountID(accountID).
	// 	SetTransferAccountID(env.Client.GetOperatorAccountID()).
	// 	FreezeWith(env.Client)
	// assert.NoError(t, err)

	// resp, err = transaction.
	// 	Sign(newKey).
	// 	Execute(env.Client)
	// assert.NoError(t, err)

	// _, err = resp.GetReceipt(env.Client)
	// assert.NoError(t, err)
}
