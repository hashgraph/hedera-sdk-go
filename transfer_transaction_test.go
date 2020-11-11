package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeTransferTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(-100)).
		SetMaxTransactionFee(HbarFrom(1, HbarUnits.Hbar)).
		SetTransactionID(testTransactionID).
		FreezeWith(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\200\302\327/\"\002\010xr\024\n\022\n\007\n\002\030\002\020\310\001\n\007\n\002\030\003\020\307\001"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"U\335!\347a\030\034\260P\376\370\323~\264jw\362\022\306}\2228\343\\\r\2506\323L\030^\201\037w\372\345\374\330\256\020;\241\006\303\270\351\350h\315\324D\255\nyH\0208\250w\007\265\252V\013">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000000transactionValidDuration:<seconds:120>cryptoTransfer:<transfers:<accountAmounts:<accountID:<accountNum:2>amount:100>accountAmounts:<accountID:<accountNum:3>amount:-100>>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

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
