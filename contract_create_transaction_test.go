package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeContractCreateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewContractCreateTransaction().
		SetAdminKey(privateKey.publicKey).
		SetInitialBalance(1e3).
		SetBytecodeFile(FileID{ File: 4 }).
		SetGas(100).
		SetProxyAccountID(AccountID{ Account: 3 }).
		SetAutoRenewPeriod(60 * 60 * 24 * 14).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionId).
		Build(mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xB7\n\002\030\004\032\"\022 \344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216 d(\350\0072\002\030\003B\000R\000Z\000"
sigMap: <
  sigPair: <
    ed25519: "\272o\034\362C\332\374\177\034\035}cG\246\360^\270s\t\234\376\301P\256\312L\233\304\240\027\234\343\3258\223\371\340z\240\306\003a\250\r\023\300\374\201_Y81\337\315\363\331,\303\310Ew\224l\004"
  >
>
transactionID: <
  transactionValidStart: <
    seconds: 124124
    nanos: 151515
  >
  accountID: <
    accountNum: 3
  >
>
nodeAccountID: <
  accountNum: 3
>
transactionFee: 1000000
transactionValidDuration: <
  seconds: 120
>
contractCreateInstance: <
  fileID: <
    fileNum: 4
  >
  adminKey: <
    ed25519: "\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"
  >
  gas: 100
  initialBalance: 1000
  proxyAccountID: <
    accountNum: 3
  >
  autoRenewPeriod: <
  >
  shardID: <
  >
  realmID: <
  >
>
`

	assert.Equal(t, txString, tx.String())
}
