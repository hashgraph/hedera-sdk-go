package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeContractDeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewContractDeleteTransaction().
		SetContractID(ContractID{Contract: 5}).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(&mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\262\001\004\n\002\030\005"
sigMap: <
  sigPair: <
    ed25519: "\002\037\252\273\3554\227\240V\217\231\347~S\204\227.\222\036\033reSJ\315?\240\224\341\272\271X\"\307\366\235\211k\360\264i<\224\313\220\343\022_\301w\201~e\376\203\227\2522|kg\202w\005"
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
contractDeleteInstance: <
  contractID: <
    contractNum: 5
  >
>
`

	assert.Equal(t, txString, tx.String())
}
