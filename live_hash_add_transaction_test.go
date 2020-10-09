package hedera

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestSerializeLiveHashAddTransaction(t *testing.T) {

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	_hash, err := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		println("Decode Failed.")
	}

	tx, err := NewLiveHashAddTransaction().
		SetAccountID(AccountID{Account: 3}).
		SetDuration((3000 * 10) * time.Millisecond).
		SetHash(_hash).
		Freeze()

	assert.NoError(t, err)

	tx.Sign(privateKey)

	//assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xz(\022\002\030\003\032\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\253\254,\271\274\307\325G;U\001\017:\264\217\224\034V\336E\320\276\035\027\315\201+0y\3125\212Kb\240Ph\263\243\372zx\251w!\257;\313<\331\204\3138\206\225\263\377Y\255T}K\020\t">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>cryptoUpdateAccount:<accountIDToUpdate:<accountNum:3>key:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestLiveHashAddTransaction_Execute(t *testing.T) {
	client, err := ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	_hash, err := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		println("Decode Failed.")
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	txID, err := NewLiveHashAddTransaction().
		SetAccountID(AccountID{Account: 3}).
		SetKeys(newKey.PublicKey()).
		SetDuration((3000 * 10) * time.Millisecond).
		SetHash(_hash).
		Execute(client)

	assert.NoError(t, err)

	println("TransactionID", txID.TransactionID.String())
	println("NodeID", txID.NodeID.String())
}
