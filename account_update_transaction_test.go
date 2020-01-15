package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeAccountUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewAccountUpdateTransaction().
		SetTransactionID(testTransactionID).
		SetAccountID(AccountID{Account: 3}).
		SetKey(privateKey.PublicKey()).
		SetMaxTransactionFee(1e6).
		Build(&mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xz(\022\002\030\003\032\"\022 \344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"
sigMap: <
  sigPair: <
    ed25519: "\253\254,\271\274\307\325G;U\001\017:\264\217\224\034V\336E\320\276\035\027\315\201+0y\3125\212Kb\240Ph\263\243\372zx\251w!\257;\313<\331\204\3138\206\225\263\377Y\255T}K\020\t"
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
cryptoUpdateAccount: <
  accountIDToUpdate: <
    accountNum: 3
  >
  key: <
    ed25519: "\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"
  >
>
`

	assert.Equal(t, txString, tx.String())
}
