package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeTokenDeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewTokenDeleteTransaction().
		SetTokenID(TokenID{Token: 3}).
		SetTransactionID(TransactionID{AccountID: AccountID{Account: 3}, ValidStart: time.Unix(0, 0)}).
		FreezeWith(mockClient)
	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\006\n\000\022\002\030\003\022\002\030\003\030\200\302\327/\"\002\010x\232\002\004\n\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\020\317\336\255\r\234XU\177^\356\221<^>\211\353E\363\013\360\233\3407\n\306\311\022&\034ih8\034I\240\347\301k\270\025\010\030\353\250\303\205\032t\233\033\354\033\321\351FO-\320\230\2569\221\000">>transactionID:<transactionValidStart:<>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000000transactionValidDuration:<seconds:120>tokenDeletion:<token:<tokenNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestTokenDeleteTransaction_Execute(t *testing.T) {
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
