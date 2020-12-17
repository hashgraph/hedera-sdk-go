package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCryptoTransferTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_CryptoTransfer_Nothing(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTransferTransaction().
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

//func Test_CryptoTransfer_1000(t *testing.T) {
//	client := newTestClient(t)
//	var err error
//	tx := make([]*TransferTransaction, 500)
//	response := make([]TransactionResponse, len(tx))
//	receipt := make([]TransactionReceipt, len(tx))
//
//	for i := 0; i < len(tx); i++ {
//		tx[i], err = NewTransferTransaction().
//			AddHbarTransfer(client.GetOperatorAccountID(), HbarFromTinybar(-10)).
//			AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(10)).
//			SetMaxTransactionFee(NewHbar(1)).
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
