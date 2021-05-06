package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCryptoTransferTransaction_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)
}

func Test_CryptoTransfer_Nothing(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)
}

func TestCryptoTransferTransaction_FlippedAmount_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(10)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(-10)).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional receipt status INVALID_SIGNATURE"), err.Error())
	}
}

func TestCryptoTransferTransaction_RepeatingAmount_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTransferTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(10)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(-10)).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(10)).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status ACCOUNT_REPEATED_IN_ACCOUNT_AMOUNTS received for transaction %s", resp.TransactionID), err.Error())
	}
}

//func Test_CryptoTransfer_1000(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//	var err error
//	tx := make([]*TransferTransaction, 500)
//	response := make([]TransactionResponse, len(tx))
//	receipt := make([]TransactionReceipt, len(tx))
//
//	for i := 0; i < len(tx); i++ {
//		tx[i], err = NewTransferTransaction().
//			AddHbarTransfer(env.Client.GetOperatorAccountID(), HbarFromTinybar(-10)).
//			AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(10)).
//			FreezeWith(env.Client)
//		if err != nil {
//			panic(err)
//		}
//
//		_, err = tx[i].SignWithOperator(env.Client)
//		if err != nil {
//			panic(err)
//		}
//
//		response[i], err = tx[i].Execute(env.Client)
//		if err != nil {
//			panic(err)
//		}
//
//		receipt[i], err = response[i].GetReceipt(env.Client)
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Printf("\r%v", i)
//	}
//}
