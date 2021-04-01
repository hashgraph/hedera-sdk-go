package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScheduleCreateTransaction_Execute(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	transactionID := TransactionIDGenerate(client.GetOperatorAccountID())

	tx := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance)

	assert.NoError(t, err)

	scheduleTx, err := tx.Schedule()
	assert.NoError(t, err)

	resp, err := scheduleTx.
		SetPayerAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewScheduleInfoQuery().
		SetScheduleID(*receipt.ScheduleID).
		SetQueryPayment(NewHbar(2)).
		Execute(client)
	assert.NoError(t, err)

	infoTx, err := info.GetTransaction()
	assert.NoError(t, err)
	assert.NotNil(t, infoTx)

	println(info.Executed.String())

	tx2, err := NewScheduleDeleteTransaction().
		SetScheduleID(*receipt.ScheduleID).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx2.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	if err != nil{
		assert.Equal(t, fmt.Sprintf("exceptional precheck status SCHEDULE_ALREADY_EXECUTED"), err.Error())
	}
}

func TestScheduleCreateTransaction_SetTransaction_Execute(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	transactionID := TransactionIDGenerate(client.GetOperatorAccountID())

	tx := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetTransactionID(transactionID).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance)

	scheduleTx, err := NewScheduleCreateTransaction().
		SetPayerAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetScheduledTransaction(tx)
	assert.NoError(t, err)

	resp, err := scheduleTx.Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	_, err = NewScheduleInfoQuery().
		SetScheduleID(*receipt.ScheduleID).
		Execute(client)
	assert.NoError(t, err)

	tx2, err := NewScheduleDeleteTransaction().
		SetScheduleID(*receipt.ScheduleID).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx2.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	if err != nil{
		assert.Equal(t, fmt.Sprintf("exceptional precheck status SCHEDULE_ALREADY_EXECUTED"), err.Error())
	}
}

func TestScheduleCreateTransaction_MultiSig_Execute(t *testing.T) {
	client := newTestClient(t, false)

	keys := make([]PrivateKey, 3)
	pubKeys := make([]PublicKey, 3)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	keyList := NewKeyList().
		AddAllPublicKeys(pubKeys)

	createResponse, err := NewAccountCreateTransaction().
		SetKey(keyList).
		SetInitialBalance(NewHbar(10)).
		Execute(client)
	assert.NoError(t, err)

	transactionReceipt, err := createResponse.GetReceipt(client)
	assert.NoError(t, err)

	transactionID := TransactionIDGenerate(client.GetOperatorAccountID())

	newAccountID := *transactionReceipt.AccountID

	transferTx := NewTransferTransaction().
		SetTransactionID(transactionID).
		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorAccountID(), HbarFrom(1, HbarUnits.Hbar))

	scheduled, err := transferTx.Schedule()
	assert.NoError(t, err)

	scheduleResponse, err := scheduled.Execute(client)
	assert.NoError(t, err)

	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	assert.NoError(t, err)

	scheduleID := *scheduleReceipt.ScheduleID

	info, err := NewScheduleInfoQuery().
		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
		SetScheduleID(scheduleID).
		Execute(client)
	assert.NoError(t, err)

	transfer, err := info.GetTransaction()
	assert.NoError(t, err)
	assert.NotNil(t, transfer)

	signTransaction, err := NewScheduleSignTransaction().
		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
		SetScheduleID(scheduleID).
		FreezeWith(client)
	assert.NoError(t, err)

	signTransaction.Sign(keys[0])
	signTransaction.Sign(keys[1])
	signTransaction.Sign(keys[2])

	resp, err := signTransaction.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info2, err := NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
		Execute(client)
	assert.NoError(t, err)
	assert.False(t, info2.Executed.IsZero())
}

func TestScheduleCreateTransaction_Delete_Execute(t *testing.T) {
	client := newTestClient(t, false)

	key, err := GeneratePrivateKey()
	key2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	createResponse, err := NewAccountCreateTransaction().
		SetKey(key).
		SetInitialBalance(NewHbar(10)).
		Execute(client)
	assert.NoError(t, err)

	transactionReceipt, err := createResponse.GetReceipt(client)
	assert.NoError(t, err)

	transactionID := TransactionIDGenerate(client.GetOperatorAccountID())

	newAccountID := *transactionReceipt.AccountID

	transferTx := NewTransferTransaction().
		SetTransactionID(transactionID).
		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
		AddHbarTransfer(client.GetOperatorAccountID(), HbarFrom(1, HbarUnits.Hbar))

	scheduled, err := transferTx.Schedule()
	assert.NoError(t, err)

	fr, err := scheduled.SetAdminKey(key2).FreezeWith(client)
	assert.NoError(t, err)

	scheduleResponse, err := fr.Sign(key2).Execute(client)
	assert.NoError(t, err)

	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
	assert.NoError(t, err)

	scheduleID := *scheduleReceipt.ScheduleID

	info, err := NewScheduleInfoQuery().
		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
		SetScheduleID(scheduleID).
		Execute(client)
	assert.NoError(t, err)

	transfer, err := info.GetTransaction()
	assert.NoError(t, err)
	assert.NotNil(t, transfer)
	assert.True(t, info.Executed.IsZero())

	tx2, err := NewScheduleDeleteTransaction().
		SetScheduleID(scheduleID).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err := tx2.
		Sign(key2).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info2, err := NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
		Execute(client)
	assert.NoError(t, err)
	assert.False(t, info2.Deleted.IsZero())
}

func TestScheduleCreateTransaction_CheckValidGetTransaction_Execute(t *testing.T) {
	client := newTestClient(t, false)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	transactionID := TransactionIDGenerate(client.GetOperatorAccountID())

	tx := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance)

	assert.NoError(t, err)

	scheduleTx, err := tx.Schedule()
	assert.NoError(t, err)

	resp, err := scheduleTx.
		SetPayerAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewScheduleInfoQuery().
		SetScheduleID(*receipt.ScheduleID).
		SetQueryPayment(NewHbar(2)).
		Execute(client)
	assert.NoError(t, err)

	infoTx, err := info.GetTransaction()
	assert.NoError(t, err)

	assert.NotNil(t, infoTx)

	switch createTx := infoTx.(type){
	case AccountCreateTransaction:
		assert.Equal(t, createTx.pbBody.GetCryptoCreateAccount().InitialBalance, uint64(NewHbar(1).tinybar))
	}

	tx2, err := NewScheduleDeleteTransaction().
		SetScheduleID(*receipt.ScheduleID).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx2.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	if err != nil{
		assert.Equal(t, fmt.Sprintf("exceptional precheck status SCHEDULE_ALREADY_EXECUTED"), err.Error())
	}
}
