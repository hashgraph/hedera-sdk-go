package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSerializeContractUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewContractUpdateTransaction().
		SetContractID(ContractID{ Contract: 3 }).
		SetAdminKey(privateKey.publicKey).
		SetBytecodeFile(FileID{ File: 5 }).
		SetExpirationTime(time.Unix(1569375111277, 0)).
		SetProxyAccountID(AccountID{ Account: 3 }).
		SetAutoRenewPeriod(60 * 60 * 24 * 14).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionId).
		Build(mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xJ;\n\002\030\003\022\007\010\355\220\257\260\326-\032\"\022 \344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\2162\002\030\003:\000B\002\030\005"
sigMap: <
  sigPair: <
    ed25519: "5\361\347\252\224\324j\371\256s.\373\3338\034\364qUV8L\004R\242\324\235\n\202p\214\342\241\324\252\260\236\037zE\353\202c\217\330\\Q\314c\323\0323\r\3653N2\373\242\202\305L\235\014\006"
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
contractUpdateInstance: <
  contractID: <
    contractNum: 3
  >
  expirationTime: <
    seconds: 1569375111277
  >
  adminKey: <
    ed25519: "\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"
  >
  proxyAccountID: <
    accountNum: 3
  >
  autoRenewPeriod: <
  >
  fileID: <
    fileNum: 5
  >
>
`

	assert.Equal(t, txString, tx.String())

}
