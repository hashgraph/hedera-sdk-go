package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeAccountCreateTransaction(t *testing.T) {
	key, err := PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	tx, err := NewAccountCreateTransaction().
		SetKey(key.PublicKey()).
		SetInitialBalance(HbarFromTinybar(450)).
		SetProxyAccountID(AccountID{Account: 1020}).
		SetReceiverSignatureRequired(true).
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		Freeze()

	assert.NoError(t, err)

	tx = tx.Sign(key)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xZI\n\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216\020\302\003\032\003\030\374\0070\377\377\377\377\377\377\377\377\1778\377\377\377\377\377\377\377\377\177@\001J\005\010\320\310\341\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\364f#\276\351\250\316\345\246\323\355\2778\2653\006\222\213P\013hx`+"`"+`\224\t\251\327\375\253\200\343K\021\235\335y\2249_\312\3574\022u\271,\331\205\010@\0172\340\252\33477\006\3029\"\007e\007">>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestAccountCreateTransaction_Execute(t *testing.T) {
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

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_AccountCreate_NoKey(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewAccountCreateTransaction().
		Execute(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status KEY_REQUIRED received for transaction %s", resp.TransactionID), err.Error())
}
