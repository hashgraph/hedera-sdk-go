package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSerializeSystemDeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewSystemDeleteTransaction().
		SetID(FileID{File: 3}).
		SetExpirationTime(time.Unix(15415151511, 0)).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(&mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\242\001\014\032\006\010\227\227\302\2669\n\002\030\003"
sigMap: <
  sigPair: <
    ed25519: "\200_Ot7\376\350\3111!\2344\257S\203i\353&\203\367\375\311]` + "`" + `f\226\202\311\275\340\\\251J\005\240\226/\251\242\351_6\210\336\001\243\305\363k'\343\314\200\002\025\337(\177\243n\336\220h\n"
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
systemDelete: <
  fileID: <
    fileNum: 3
  >
  expirationTime: <
    seconds: 15415151511
  >
>
`

	assert.Equal(t, txString, tx.String())
}
