package hedera

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//func TestSerializeLiveHashAddTransaction(t *testing.T) {

//	privateKey, err := PrivateKeyFromString(mockPrivateKey)
//	assert.NoError(t, err)

//	client, err := ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

//	if err != nil {
//		client = ClientForTestnet()
//	}

//	_hash, err := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")
//	if err != nil {

//	}

//	configOperatorID := os.Getenv("OPERATOR_ID")
//	configOperatorKey := os.Getenv("OPERATOR_KEY")

//	if configOperatorID != "" && configOperatorKey != "" {
//		operatorAccountID, err := AccountIDFromString(configOperatorID)
//		assert.NoError(t, err)

//		operatorKey, err := PrivateKeyFromString(configOperatorKey)
//		assert.NoError(t, err)

//		client.SetOperator(operatorAccountID, operatorKey)
//	}

//	if err != nil {

//	}

//	tx, err := NewLiveHashAddTransaction().
//		SetAccountID(AccountID{Account: 3}).
//		SetDuration((3000 * 10) * time.Millisecond).
//		SetHash(_hash).
//		FreezeWith(client)

//	assert.NoError(t, err)

//	tx.Sign(privateKey)

//	//assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xz(\022\002\030\003\032\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\253\254,\271\274\307\325G;U\001\017:\264\217\224\034V\336E\320\276\035\027\315\201+0y\3125\212Kb\240Ph\263\243\372zx\251w!\257;\313<\331\204\3138\206\225\263\377Y\255T}K\020\t">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>cryptoUpdateAccount:<accountIDToUpdate:<accountNum:3>key:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
//}

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
