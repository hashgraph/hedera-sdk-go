package hedera

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeFileCreateTransaction(t *testing.T) {
	date := time.Unix(1554158542, 0)

	key, err := Ed25519PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	tx := NewFileCreateTransaction().
		AddKey(key.PublicKey()).
		SetContents(Bytes([]byte{1, 2, 3, 4})).
		SetExpirationTime(date).
		SetNodeAccountID(AccountID{Account: 3}).
		SetTransactionID(TransactionID{
			AccountID:  AccountID{Account: 2},
			ValidStart: date,
		}).
		SetMaxTransactionFee(100_000).
		Build(nil)

	assert.NoError(t, err)

	tx.Sign(key)

	txString := `bodyBytes: "\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\030\240\215\006\"\002\010x\212\0014\022\006\010\316\247\212\345\005\032$\n\"\022 \344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216\"\004\001\002\003\004"
sigMap: <
  sigPair: <
    ed25519: "\347\020\300\212\213D\306\370\205Al\375\305\347\235\313\276by\237\230\270\226)Z\370\362\271\353W&\035\251\033u\376\227\321?5\263\355f\203\250\304\315~\317\312\272\352P\316\331\355\031\256\006t\374B_\006"
  >
>
transactionID: <
  transactionValidStart: <
    seconds: 1554158542
  >
  accountID: <
    accountNum: 2
  >
>
nodeAccountID: <
  accountNum: 3
>
transactionFee: 100000
transactionValidDuration: <
  seconds: 120
>
fileCreate: <
  expirationTime: <
    seconds: 1554158542
  >
  keys: <
    keys: <
      ed25519: "\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"
    >
  >
  contents: "\001\002\003\004"
>
`

	assert.Equal(t, txString, tx.String())
}
