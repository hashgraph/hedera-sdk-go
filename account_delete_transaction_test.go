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
		SetAccountId(AccountID{ Account: 3 }).
		// todo: transfer account id
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionId).
		Build(mockClient)

	tx.Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xb\004\022\002\030\003"
sigMap: <
  sigPair: <
    ed25519: "IZ\225\325\352\274\200K~\207\321>\204\034\003\231O\254\2435\324\257a\033\272\247\031\231\240R\303\375\356m\373\320\235\312\353\363\027\2411'\362*\236#\215\362\372;f\342\332\257=W\333n\213\017*\002"
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
  deleteAccountID: <
    accountNum: 3
  >
>
`

	assert.Equal(t, txString, tx.String())
}
