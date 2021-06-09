package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}

func Test_ReceiptTransaction_InvalidTransactionID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	key, err := GeneratePrivateKey()
	assert.NoError(t, err)

	response, err := NewAccountCreateTransaction().
		SetKey(key).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	tx := NewTransferTransaction().
		AddHbarTransfer(accountID, NewHbar(1).Negated()).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(1))

	scheduleTx, err := tx.Schedule()
	assert.NoError(t, err)

	scheduleTx = scheduleTx.
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetPayerAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID()))

	response, err = scheduleTx.Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	assert.NoError(t, err)

	scheduleID := *receipt.ScheduleID

	scheduleSignTx, err := NewScheduleSignTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	response, err = scheduleSignTx.Execute(env.Client)

	_, err = response.GetReceipt(env.Client)
	assert.Error(t, err)

	switch receiptErr := err.(type) {
	case ErrHederaReceiptStatus:
		assert.NotNil(t, receiptErr.Receipt.ExchangeRate)
	default:
		panic("err was not a `ErrHederaReceiptStatus")
	}
}
