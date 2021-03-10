package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCryptoTransferTransaction_Execute(t *testing.T) {
	client := newTestClient(t, false)

	resp, err := NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_CryptoTransfer_Nothing(t *testing.T) {
	client := newTestClient(t, false)

	resp, err := NewTransferTransaction().
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func TestCryptoTransferTransaction_FlippedAmount_Execute(t *testing.T) {
	client := newTestClient(t, false)

	resp, err := NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(10)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(-10)).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_SIGNATURE"), err.Error())
	}
}

func TestCryptoTransferTransaction_RepeatingAmount_Execute(t *testing.T) {
	client := newTestClient(t, false)

	resp, err := NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(10)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(-10)).
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(10)).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_REPEATED_IN_ACCOUNT_AMOUNTS received for transaction %s", resp.TransactionID), err.Error())
	}
}

//func Test_CryptoTransfer_1000(t *testing.T) {
//	client := newTestClient(t, false)
//	var err error
//	tx := make([]*TransferTransaction, 500)
//	response := make([]TransactionResponse, len(tx))
//	receipt := make([]TransactionReceipt, len(tx))
//
//	for i := 0; i < len(tx); i++ {
//		tx[i], err = NewTransferTransaction().
//			AddHbarTransfer(client.GetOperatorAccountID(), HbarFromTinybar(-10)).
//			AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(10)).
//			FreezeWith(client)
//		if err != nil {
//			panic(err)
//		}
//
//		_, err = tx[i].SignWithOperator(client)
//		if err != nil {
//			panic(err)
//		}
//
//		response[i], err = tx[i].Execute(client)
//		if err != nil {
//			panic(err)
//		}
//
//		receipt[i], err = response[i].GetReceipt(client)
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Printf("\r%v", i)
//	}
//}
