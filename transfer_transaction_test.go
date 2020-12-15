package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCryptoTransferTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(-1)).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func TestCryptoTransferTransactionNothing_Execute(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTransferTransaction().
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func TestCryptoTransferTransaction1000_Execute(t *testing.T) {
	client := newTestClient(t)
	var err error
	tx := make([]*TransferTransaction, 1000)
	response := make([]TransactionResponse, len(tx))
	receipt := make([]TransactionReceipt, len(tx))

	for i := 0; i < len(tx); i++ {
		tx[i], err = NewTransferTransaction().
			AddHbarTransfer(client.GetOperatorAccountID(), HbarFromTinybar(10)).
			AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(-10)).
			SetMaxTransactionFee(NewHbar(1)).
			FreezeWith(client)
		assert.NoError(t, err)

		tx[i], err = tx[i].SignWithOperator(client)
		assert.NoError(t, err)
		response[i], err = tx[i].Execute(client)
		assert.NoError(t, err)
		receipt[i], err = response[i].GetReceipt(client)
		assert.NoError(t, err)
	}

	//for i := 0; i < len(tx); i++ {
	//	response[i], err = tx[i].Execute(client)
	//	assert.NoError(t, err)
	//}

	//for i := 0; i < len(tx); i++ {
	//	receipt[i], err = response[i].GetReceipt(client)
	//	assert.NoError(t, err)
	//}

}
