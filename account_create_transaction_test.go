package hedera

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeAccountCreateTransaction(t *testing.T) {
	date := time.Unix(1554158542, 0)

	key, err := Ed25519PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	tx := NewAccountCreateTransaction().
		SetKey(key.PublicKey()).
		SetInitialBalance(450).
		// SetProxyAccountID(AccountID{Account: 1020}).
		// SetReceiverSignatureRequired(true).
		SetNodeAccountID(AccountID{Account: 3}).
		SetTransactionID(TransactionID{
			AccountID:  AccountID{Account: 2},
			ValidStart: date,
		}).
		SetMaxTransactionFee(100_000).
		Build(nil)

	assert.NoError(t, err)

	tx.Sign(key)

	txString := `bodyBytes: "\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\030\240\215\006\"\002\010xZB\n\"\022 \344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216\020\302\0030\377\377\377\377\377\377\377\377\1778\377\377\377\377\377\377\377\377\177J\005\010\320\310\341\003"
sigMap: <
  sigPair: <
    ed25519: "\362<\255\304\241>\035\2775J\306w\377\033k\031\217\271i\263\301\357z!H\037\016Yp31HM\033W\317\303\317\247W\233\003\030\330&\362*9\346x\227\211r\272\"t\310\373KB\016\242\275\001"
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
cryptoCreateAccount: <
  key: <
    ed25519: "\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"
  >
  initialBalance: 450
  sendRecordThreshold: 9223372036854775807
  receiveRecordThreshold: 9223372036854775807
  autoRenewPeriod: <
    seconds: 7890000
  >
>
`

	assert.Equal(t, txString, tx.String())
}
