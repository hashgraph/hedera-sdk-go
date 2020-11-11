package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeTokenCreateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewTokenCreateTransaction().
		SetTreasuryAccountID(AccountID{Account: 3}).
		SetMaxTransactionFee(NewHbar(1000)).
		SetTransactionID(TransactionID{AccountID: AccountID{Account: 3}, ValidStart: time.Unix(0, 0)}).
		SetExpirationTime(time.Unix(0, 0)).
		SetNodeAccountIDs([]AccountID{{0, 0, 3}}).
		FreezeWith(mockClient)
	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\006\n\000\022\002\030\003\022\002\030\003\030\200\320\333\303\364\002\"\002\010x\352\001\t*\002\030\003x\320\310\341\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\376\256~\275\264\200<jB\254\213\243On1\210\372Q\340\037\360\016z\223\315pm\330\251qC*)\301\261\373\363\205h\247\\\327\346\363\304\037^\3659\026-\260\372\255/\"\211%\363t\376\202\374\017">>transactionID:<transactionValidStart:<>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000000000transactionValidDuration:<seconds:120>tokenCreation:<treasury:<accountNum:3>autoRenewPeriod:7890000>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestTokenCreateTransaction_Execute(t *testing.T) {
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

	resp, err = NewTokenDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
