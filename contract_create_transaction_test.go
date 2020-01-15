package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSerializeContractCreateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewContractCreateTransaction().
		SetAdminKey(privateKey.PublicKey()).
		SetInitialBalance(1e3).
		SetBytecodeFile(FileID{File: 4}).
		SetGas(100).
		SetProxyAccountID(AccountID{Account: 3}).
		SetAutoRenewPeriod(60 * 60 * 24 * 14 * time.Second).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(&mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xB;\n\002\030\004\032\"\022 \344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216 d(\350\0072\002\030\003B\004\010\200\352IR\000Z\000"
sigMap: <
  sigPair: <
    ed25519: "\200\207\214\233\270\235$\302\347\343\002\020\333&\343\363_!@>2\354\351\311.)\223\345\230\347:\241\216,\371\262\300\332U2\341\352\337B\311\216t>\233\023\217\366\337\242i\264\354+\026,P\001\363\001"
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
    seconds: 1209600
  >
  shardID: <
  >
  realmID: <
  >
>
`

	assert.Equal(t, txString, tx.String())
}
