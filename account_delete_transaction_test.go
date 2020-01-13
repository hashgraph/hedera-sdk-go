package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeAccountDeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewAccountDeleteTransaction().
		SetDeleteAccountId(AccountID{Account: 3}).
		SetTransferAccountID(AccountID{Account: 2}).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	tx.Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xb\010\n\002\030\002\022\002\030\003"
sigMap: <
  sigPair: <
    ed25519: "&\321\261A\177f\316\346\326\346\t\004\202\272\365Q_/\027\014-:\3429eM\265\263\275N\227\350?G\270f\347\205mk0\211zH\3244w\221\213\005\315\1776\236~Z\341\2138\277TLF\007"
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
cryptoDelete: <
  transferAccountID: <
    accountNum: 2
  >
  deleteAccountID: <
    accountNum: 3
  >
>
`

	assert.Equal(t, txString, tx.String())
}
