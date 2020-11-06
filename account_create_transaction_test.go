package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeAccountCreateTransaction(t *testing.T) {
	date := time.Unix(1554158542, 0)

	key, err := PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	id, err := AccountIDFromString("3")

	tx, err := NewAccountCreateTransaction().
		SetKey(key.PublicKey()).
		SetInitialBalance(HbarFromTinybar(450)).
		SetProxyAccountID(AccountID{Account: 1020}).
		SetReceiverSignatureRequired(true).
		SetNodeAccountIDs([]AccountID{id}).
		SetTransactionID(TransactionID{
			AccountID:  AccountID{Account: 2},
			ValidStart: date,
		}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		Freeze()

	assert.NoError(t, err)

	tx = tx.Sign(key)

	assert.Equal(t, `bodyBytes:"\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\000\030\300\204=\"\002\010xZI\n\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216\020\302\003\032\003\030\374\0070\377\377\377\377\377\377\377\377\1778\377\377\377\377\377\377\377\377\177@\001J\005\010\320\310\341\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\327\233\322+(\314\017\301T\234~S\312\207\361X\355\252{2\313\332\0248\367\334W\010\260;\223\t\313O\\_:\213\362\224\377\205[d)u%A\216S\235~^XW2\214\354\363\346*j\273\007">>transactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>nodeAccountID:<>transactionFee:1000000transactionValidDuration:<seconds:120>cryptoCreateAccount:<key:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">initialBalance:450proxyAccountID:<accountNum:1020>sendRecordThreshold:9223372036854775807receiveRecordThreshold:9223372036854775807receiverSigRequired:trueautoRenewPeriod:<seconds:7890000>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
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

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs(nodeIDs).
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
