package hedera

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeLiveHashAddTransaction(t *testing.T) {
	client, err := newMockClient()
	assert.NoError(t, err)

	newKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewLiveHashAddTransaction().
		SetAccountID(AccountID{Account: 3}).
		SetDuration((3000 * 10) * time.Millisecond).
		SetTransactionID(TransactionID{AccountID: AccountID{Account: 3}, ValidStart: time.Unix(0,0)}).
		FreezeWith(client)

	assert.NoError(t, err)

	tx.Sign(newKey)

	assert.Equal(t, `bodyBytes:"\n\006\n\000\022\002\030\003\022\002\030\003\030\200\302\327/\"\002\010xR\n\032\010\n\002\030\003*\002\010\036"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\215\276Id\212\326\363f\353q\315l\211\207\334\245\244\203E\031\276\343q\236\203f\211\210c/\334\224\213+f\336\200\025{7\007\331\246/\206r\211\r\305J\212\3470\232\271G\301\2271\346\025\005">>transactionID:<transactionValidStart:<>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000000transactionValidDuration:<seconds:120>cryptoAddLiveHash:<liveHash:<accountId:<accountNum:3>duration:<seconds:30>>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestLiveHashAddTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	_hash, err := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {

	}

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	resp, err = NewLiveHashAddTransaction().
		SetAccountID(accountID).
		SetDuration(24 * 30 * time.Hour).
		SetNodeAccountIDs(nodeIDs).
		SetHash(_hash).
		SetKeys(newKey.PublicKey()).
		Execute(client)

	assert.Error(t, err)

	resp, err = NewLiveHashDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetHash(_hash).
		Execute(client)
	assert.Error(t, err)

	resp, err = NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetTransferAccountID(client.GetOperatorAccountID()).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
