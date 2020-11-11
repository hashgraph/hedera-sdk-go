package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeAccountDeleteTransaction(t *testing.T) {
	mockClient := newTestClient(t)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(AccountID{Account: 3}).
		SetTransferAccountID(AccountID{Account: 2}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetNodeAccountIDs([]AccountID{{0, 0, 6}}).
		FreezeWith(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\006\030\300\204=\"\002\010xb\010\n\002\030\002\022\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\251\320L\230\362|]\3504&\300\220\341Y\206<\002\003\270Ut\376\366\203\033\35116.\251y\275\2269\002v\222\033\3264\273\235K\243\216\333\201\261?\211\022fT\330\351/\020\342g\352dM\346\t">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:6>transactionFee:1000000transactionValidDuration:<seconds:120>cryptoDelete:<transferAccountID:<accountNum:2>deleteAccountID:<accountNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestAccountDeleteTransaction_Execute(t *testing.T) {
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

	accountID := receipt.AccountID
	assert.NoError(t, err)

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(*accountID).
		SetTransferAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs(nodeIDs).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(TransactionIDGenerate(*accountID)).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = tx.Sign(newKey).Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
