package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScheduleCreateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Build(client)
	assert.NoError(t, err)

	scheduleTx := tx.Schedule()

	frozen, err := scheduleTx.
		SetPayerAccountID(client.GetOperatorID()).
		SetAdminKey(newKey.PublicKey()).
		Build(client)
	assert.NoError(t, err)

	frozen = frozen.Sign(newKey)

	txID, err := frozen.Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	_, err = NewScheduleInfoQuery().
		SetScheduleID(receipt.GetScheduleID()).
		SetMaxQueryPayment(NewHbar(2)).
		Execute(client)
	assert.NoError(t, err)

	tx2, err := NewScheduleDeleteTransaction().
		SetScheduleID(receipt.GetScheduleID()).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx2.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}

func TestScheduleCreateTransaction_SetTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		SetNodeAccountID(AccountID{0,0,3}).
		Build(client)
	assert.NoError(t, err)

	scheduleTx, err := NewScheduleCreateTransaction().
		SetTransaction(tx).
		SetAdminKey(newKey.PublicKey()).
		SetPayerAccountID(client.GetOperatorID()).
		Build(client)
	assert.NoError(t, err)

	scheduleTx = scheduleTx.Sign(newKey)

	txID, err := scheduleTx.Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	_, err = NewScheduleInfoQuery().
		SetScheduleID(receipt.GetScheduleID()).
		Execute(client)
	assert.NoError(t, err)

	tx2, err := NewScheduleDeleteTransaction().
		SetScheduleID(receipt.GetScheduleID()).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx2.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}

func TestScheduleCreateTransaction_Signature_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		SetNodeAccountID(AccountID{Account: 3}).
		Build(client)
	assert.NoError(t, err)

	txBytes, err := tx.MarshalBinary()
	assert.NoError(t, err)

	signature1, err := newKey.SignTransaction(&tx)

	err = tx.UnmarshalBinary(txBytes)
	assert.NoError(t, err)

	tx2 := tx.Schedule()

	tx3, err := tx2.
		SetAdminKey(newKey.PublicKey()).
		SetPayerAccountID(client.GetOperatorID()).
		SetTransactionValidDuration(30 * time.Second).
		Build(client)
	assert.NoError(t, err)

	tx3 = tx3.Sign(newKey)

	resp, err := tx3.Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	scheduleID := receipt.GetScheduleID()

	resp, err = NewScheduleSignTransaction().
		SetScheduleID(scheduleID).
		SetNodeAccountID(AccountID{Account: 3}).
		AddScheduleSignature(newKey.PublicKey(), signature1).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tx4, err := NewScheduleDeleteTransaction().
		SetScheduleID(scheduleID).
		Build(client)
	assert.NoError(t, err)

	resp, err = tx4.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
