package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeTokenTransferTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewTransferTransaction().
		AddTokenTransfer(TokenID{Token: 3}, AccountID{Account: 3}, -10).
		AddTokenTransfer(TokenID{Token: 3}, AccountID{Account: 4}, 10).
		SetTransactionID(TransactionID{AccountID: AccountID{Account: 3}, ValidStart: time.Unix(0, 0)}).
		FreezeWith(mockClient)
	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\006\n\000\022\002\030\003\022\002\030\003\030\200\302\327/\"\002\010xr\030\n\000\022\024\n\002\030\003\022\006\n\002\030\003\020\023\022\006\n\002\030\004\020\024"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\037\217'Q\271S#\034\031R\020\241P\343\025+\270\014h'4oe^\344`+"`"+`d\256\341o\214u\211\356\257\340,\335C\177\361C\222\323/\352\255\203\010\377\246{\022Dy\031\244\230\251\001\210\000\306\017">>transactionID:<transactionValidStart:<>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000000transactionValidDuration:<seconds:120>cryptoTransfer:<transfers:<>tokenTransfers:<token:<tokenNum:3>transfers:<accountID:<accountNum:3>amount:-10>transfers:<accountID:<accountNum:4>amount:10>>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestTokenTransferTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GeneratePrivateKey()
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

	resp, err = NewTokenCreateTransaction().
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

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	nodeId := resp.NodeID

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		AddTokenTransfer(tokenID, client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID, accountID, 10).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewTokenWipeTransaction().
		SetNodeAccountIDs([]AccountID{nodeId}).
		SetTokenID(tokenID).
		SetAccountID(accountID).
		SetAmount(10).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
