package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeCryptoTransferTransaction(t *testing.T) {
	tx, err := newMockTransaction()

	assert.NoError(t, err)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xr\024\n\022\n\007\n\002\030\002\020\307\001\n\007\n\002\030\003\020\310\001"
sigMap: <
  sigPair: <
    ed25519: "\"\352\272~\022\246\274\024\330\364v\036\376N\217\350X\253\370\377\324e\317 \312\262\341\353a+}\341\350]!\264\033\010\234Z\246\311\204\030\215'\0338I\350{\334\251\355\225@\255|X\340U\210\260\007"
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
cryptoTransfer: <
  transfers: <
    accountAmounts: <
      accountID: <
        accountNum: 2
      >
      amount: -100
    >
    accountAmounts: <
      accountID: <
        accountNum: 3
      >
      amount: 100
    >
  >
>
`

	assert.Equal(t, txString, tx.String())
}
