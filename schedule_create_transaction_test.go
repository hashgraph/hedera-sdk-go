package hedera

//import (
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestScheduleCreateTransaction_Execute(t *testing.T) {
//	client := newTestClient(t)
//
//	newKey, err := GenerateEd25519PrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(1)
//
//	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	scheduleTx, err := NewAccountCreateTransaction().
//		SetKey(newKey.PublicKey()).
//		SetMaxTransactionFee(NewHbar(2)).
//		SetInitialBalance(newBalance).
//		Schedule()
//	assert.NoError(t, err)
//
//	frozen, err := scheduleTx.
//		SetPayerAccountID(client.GetOperatorID()).
//		SetAdminKey(newKey.PublicKey()).
//		Build(client)
//	assert.NoError(t, err)
//
//	frozen = frozen.Sign(newKey)
//
//	txID, err := frozen.Execute(client)
//	assert.NoError(t, err)
//
//	receipt, err := txID.GetReceipt(client)
//	assert.NoError(t, err)
//
//	info, err := NewScheduleInfoQuery().
//		SetScheduleID(receipt.GetScheduleID()).
//		SetMaxQueryPayment(NewHbar(2)).
//		Execute(client)
//	assert.NoError(t, err)
//
//	infoTx, err := info.GetScheduledTransaction()
//	assert.NoError(t, err)
//	assert.NotNil(t, infoTx)
//
//	assert.False(t, info.ExecutedAt.IsZero())
//}
//
//func TestScheduleCreateTransaction_SetTransaction_Execute(t *testing.T) {
//	client := newTestClient(t)
//
//	newKey, err := GenerateEd25519PrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(1)
//
//	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	tx := NewAccountCreateTransaction().
//		SetKey(newKey.PublicKey()).
//		SetMaxTransactionFee(NewHbar(2)).
//		SetInitialBalance(newBalance).
//		SetNodeAccountID(AccountID{0, 0, 3})
//
//	scheduleTx, err := NewScheduleCreateTransaction().
//		SetScheduledTransaction(&tx)
//	assert.NoError(t, err)
//
//	txS, err := scheduleTx.
//		SetAdminKey(newKey.PublicKey()).
//		SetPayerAccountID(client.GetOperatorID()).
//		Build(client)
//	assert.NoError(t, err)
//
//	txS = txS.Sign(newKey)
//
//	txID, err := txS.Execute(client)
//	assert.NoError(t, err)
//
//	receipt, err := txID.GetReceipt(client)
//	assert.NoError(t, err)
//
//	info, err := NewScheduleInfoQuery().
//		SetScheduleID(receipt.GetScheduleID()).
//		Execute(client)
//	assert.NoError(t, err)
//
//	infoTx, err := info.GetScheduledTransaction()
//	assert.NoError(t, err)
//	assert.NotNil(t, infoTx)
//
//	assert.False(t, info.ExecutedAt.IsZero())
//
//}
//
//func TestScheduleCreateTransaction_MultiSig_Execute(t *testing.T) {
//	client := newTestClient(t)
//
//	keys := make([]Ed25519PrivateKey, 3)
//	pubKeys := make([]PublicKey, 3)
//
//	for i := range keys {
//		newKey, err := GenerateEd25519PrivateKey()
//		assert.NoError(t, err)
//
//		keys[i] = newKey
//		pubKeys[i] = newKey.PublicKey()
//	}
//
//	keyList := NewKeyList().
//		AddAllPublicKeys(pubKeys)
//
//	createResponse, err := NewAccountCreateTransaction().
//		SetKey(keyList).
//		SetInitialBalance(NewHbar(10)).
//		Execute(client)
//	assert.NoError(t, err)
//
//	transactionReceipt, err := createResponse.GetReceipt(client)
//	assert.NoError(t, err)
//
//	transactionID := NewTransactionID(client.GetOperatorID())
//
//	newAccountID := transactionReceipt.GetAccountID()
//
//	transferTx := NewTransferTransaction().
//		SetTransactionID(transactionID).
//		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
//		AddHbarTransfer(client.GetOperatorID(), HbarFrom(1, HbarUnits.Hbar))
//
//	scheduled, err := transferTx.Schedule()
//	assert.NoError(t, err)
//
//	scheduleResponse, err := scheduled.Execute(client)
//	assert.NoError(t, err)
//
//	scheduleReceipt, err := scheduleResponse.GetReceipt(client)
//	assert.NoError(t, err)
//
//	scheduleID := scheduleReceipt.GetScheduleID()
//
//	info, err := NewScheduleInfoQuery().
//		SetScheduleID(scheduleID).
//		Execute(client)
//	assert.NoError(t, err)
//
//	transfer, err := info.GetScheduledTransaction()
//	assert.NoError(t, err)
//	assert.NotNil(t, transfer)
//
//	signTransaction, err := NewScheduleSignTransaction().
//		SetScheduleID(scheduleID).
//		Build(client)
//	assert.NoError(t, err)
//
//	signTransaction.Sign(keys[0])
//	signTransaction.Sign(keys[1])
//	signTransaction.Sign(keys[2])
//
//	resp, err := signTransaction.Execute(client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(client)
//	assert.NoError(t, err)
//
//	info2, err := NewScheduleInfoQuery().
//		SetScheduleID(scheduleID).
//		Execute(client)
//	assert.NoError(t, err)
//	assert.False(t, info2.ExecutedAt.IsZero())
//}
