package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestSerializeAccountUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetTransactionID(testTransactionID).
		SetAccountID(AccountID{Account: 3}).
		SetKey(privateKey.PublicKey()).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xz(\022\002\030\003\032\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\253\254,\271\274\307\325G;U\001\017:\264\217\224\034V\336E\320\276\035\027\315\201+0y\3125\212Kb\240Ph\263\243\372zx\251w!\257;\313<\331\204\3138\206\225\263\377Y\255T}K\020\t">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>cryptoUpdateAccount:<accountIDToUpdate:<accountNum:3>key:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestAccountUpdateTransaction_Execute(t *testing.T) {
	client, err := ClientFromFile(os.Getenv("CONFIG_FILE"))

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

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newKey2, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	txID, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	accountID := receipt.GetAccountID()
	assert.NoError(t, err)

	tx, err := NewAccountUpdateTransaction().
		SetAccountID(accountID).
		SetKey(newKey2.PublicKey()).
		SetMaxTransactionFee(NewHbar(1)).
		Build(client)
	assert.NoError(t, err)

	tx.Sign(newKey)
	tx.Sign(newKey2)

	txID, err = tx.Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newKey2.PublicKey(), info.Key)

	tx, err = NewAccountDeleteTransaction().
		SetDeleteAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetMaxTransactionFee(NewHbar(1)).
		Build(client)
	assert.NoError(t, err)

	tx.Sign(newKey2)

	txID, err = tx.Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)

	assert.NoError(t, err)
}
