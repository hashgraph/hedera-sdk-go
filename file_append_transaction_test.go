package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestSerializeFileAppendTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewFileAppendTransaction().
		SetFileID(FileID{File: 5}).
		SetContents([]byte("This is some random data")).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\202\001\036\022\002\030\005\"\030Thisissomerandomdata"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\211Q\013\365\342\\\257B\255\370\347t\234`+"`"+`)\271\017\\\273\266\367\347\214\256]D\2004\220kC$:\252\245\227\257\351\365\344\236\244\032\336@\263a\353\001\276\257\300)x\254\021\032\217\223DF\316\016\004">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>fileAppend:<fileID:<fileNum:5>contents:"Thisissomerandomdata">`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestFileAppendTransaction_Execute(t *testing.T) {
	client, err := ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := Ed25519PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	client.SetMaxTransactionFee(NewHbar(2))

	txID, err := NewFileCreateTransaction().
		AddKey(client.GetOperatorKey()).
		SetContents([]byte("Hello")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.fileID
	assert.NotNil(t, fileID)

	txID, err = NewFileAppendTransaction().
		SetFileID(*fileID).
		SetContents([]byte(" world!")).
		Execute(client)

	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	contents, err := NewFileContentsQuery().
		SetFileID(*fileID).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, []byte("Hello world!"), contents)

	txID, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
