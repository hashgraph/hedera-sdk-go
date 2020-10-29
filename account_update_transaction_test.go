package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// func TestSerializeAccountUpdateTransaction(t *testing.T) {
// 	mockClient, err := newMockClient()
// 	assert.NoError(t, err)

// 	privateKey, err := PrivateKeyFromString(mockPrivateKey)
// 	assert.NoError(t, err)

// 	tx, err := NewAccountUpdateTransaction().
// 		SetAccountID(AccountID{Account: 3}).
// 		SetKey(privateKey.PublicKey()).
// 		SetNodeAccountID(AccountID{Account: 5}).
// 		SetMaxTransactionFee(HbarFromTinybar(1e6)).
// 		SetTransactionID(testTransactionID).
// 		FreezeWith(mockClient)

// 	assert.NoError(t, err)

// 	tx.Sign(privateKey)

// 	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xz(\022\002\030\003\032\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\253\254,\271\274\307\325G;U\001\017:\264\217\224\034V\336E\320\276\035\027\315\201+0y\3125\212Kb\240Ph\263\243\372zx\251w!\257;\313<\331\204\3138\206\225\263\377Y\255T}K\020\t">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>cryptoUpdateAccount:<accountIDToUpdate:<accountNum:3>key:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
// }

func TestAccountUpdateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newKey2, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetKey(newKey2.PublicKey()).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)

	assert.NoError(t, err)

	tx.Sign(newKey)
	tx.Sign(newKey2)

	resp, err = tx.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey(), info.Key)

	txDelete, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetNodeAccountIDs(nodeIDs).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)

	assert.NoError(t, err)

	txDelete.Sign(newKey2)

	resp, err = txDelete.Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)

	assert.NoError(t, err)
}
