package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeTokenUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewTokenUpdateTransaction().
		SetTokenID(TokenID{Token: 3}).
		SetTokenSymbol("A").
		SetTransactionID(TransactionID{AccountID: AccountID{Account: 3}, ValidStart: time.Unix(0, 0)}).
		FreezeWith(mockClient)
	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\006\n\000\022\002\030\003\022\002\030\003\030\200\302\327/\"\002\010x\242\002\007\n\002\030\003\022\001A"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\216\242\276x\346[\020K_/\220\257[\022\216\223\350\343k\217T\031\234\210\002\002\247\025\321]\347\216/\336,\340\214\305\363=\361\330\303\250\357\260\314!\223\330K\301\027\214\206\221E\221\331\"\266\"\325\016">>transactionID:<transactionValidStart:<>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000000transactionValidDuration:<seconds:120>tokenUpdate:<token:<tokenNum:3>symbol:"A">`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestTokenUpdateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetFreezeKey(client.GetOperatorPublicKey()).
		SetWipeKey(client.GetOperatorPublicKey()).
		SetKycKey(client.GetOperatorPublicKey()).
		SetSupplyKey(client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetMaxTransactionFee(NewHbar(1000)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
