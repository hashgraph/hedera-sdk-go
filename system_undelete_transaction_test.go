package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeSystemUndeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewSystemUndeleteTransaction().
		SetID(FileID{File: 3}).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(&mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\252\001\004\n\002\030\003"
sigMap: <
  sigPair: <
    ed25519: "R\254\361'\030i\300/\341I-y\347}\263\247,u_e\010\301\201\324s;\232\305}8h\3573\225\213\227H\0315_\210b` + "`" + `B0C\236kX\304\365\211\216E\221\341\311\034t\004-.\031\t"
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
systemUndelete: <
  fileID: <
    fileNum: 3
  >
>
`

	assert.Equal(t, txString, tx.String())
}
