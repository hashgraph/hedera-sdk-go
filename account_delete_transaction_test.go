package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeAccountDeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetDeleteAccountID(AccountID{Account: 3}).
		SetTransferAccountID(AccountID{Account: 2}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xb\010\n\002\030\002\022\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"&\321\261A\177f\316\346\326\346\t\004\202\272\365Q_/\027\014-:\3429eM\265\263\275N\227\350?G\270f\347\205mk0\211zH\3244w\221\213\005\315\1776\236~Z\341\2138\277TLF\007">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>cryptoDelete:<transferAccountID:<accountNum:2>deleteAccountID:<accountNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestAccountDeleteTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	txID, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	accountID := receipt.GetAccountID()
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetDeleteAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(NewTransactionID(accountID)).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
