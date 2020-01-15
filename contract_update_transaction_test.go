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
		SetContractID(ContractID{Contract: 3}).
		SetAdminKey(privateKey.PublicKey()).
		SetBytecodeFile(FileID{File: 5}).
		SetExpirationTime(time.Unix(1569375111277, 0)).
		SetProxyAccountID(AccountID{Account: 3}).
		SetAutoRenewPeriod(60 * 60 * 24 * 14 * time.Second).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(&mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xJ?\n\002\030\003\022\007\010\355\220\257\260\326-\032\"\022 \344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\2162\002\030\003:\004\010\200\352IB\002\030\005"
sigMap: <
  sigPair: <
    ed25519: "\023i\223\311Iem\320^\202}I\2127\327\270'\240#\354\206X\2053\340!\2010#j\275r\335s\216O\371\226\211\247\242\356\004\374L)\323\352\251\252C+\303L\346~k\222G\216\212\325L\n"
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
    seconds: 1209600
  >
  fileID: <
    fileNum: 5
  >
>
`

	assert.Equal(t, txString, tx.String())

}
