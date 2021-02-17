package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

	frozen, err := scheduleTx.SetAdminKey(newKey.PublicKey()).Build(client)
	assert.NoError(t, err)

	frozen = frozen.Sign(newKey)

	txID, err := frozen.Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	println(receipt.GetScheduleID().String())

	_, err = NewScheduleInfoQuery().
		SetScheduleID(receipt.GetScheduleID()).
		SetMaxQueryPayment(NewHbar(2)).
		Execute(client)
	assert.NoError(t, err)

	//println(info.)

	tx2, err := NewScheduleDeleteTransaction().
		SetScheduleID(receipt.GetScheduleID()).
		Build(client)
	assert.NoError(t, err)

	println(tx2.body().String())

	txID, err = tx2.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
