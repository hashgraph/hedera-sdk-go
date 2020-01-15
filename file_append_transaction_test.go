package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeFileAppendTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewFileAppendTransaction().
		SetFileID(FileID{File: 5}).
		SetContents([]byte("This is some random data")).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\202\001\036\022\002\030\005\"\030This is some random data"
sigMap: <
  sigPair: <
    ed25519: "\211Q\013\365\342\\\257B\255\370\347t\234` + "`" + `)\271\017\\\273\266\367\347\214\256]D\2004\220kC$:\252\245\227\257\351\365\344\236\244\032\336@\263a\353\001\276\257\300)x\254\021\032\217\223DF\316\016\004"
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
fileAppend: <
  fileID: <
    fileNum: 5
  >
  contents: "This is some random data"
>
`

	assert.Equal(t, txString, tx.String())
}
