package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeContractExecuteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewContractExecuteTransaction().
		SetContractID(ContractID{Contract: 5}).
		SetGas(141).
		SetPayableAmount(HbarFromTinybar(10000)).
		// this was pulled from the Java test
		SetFunctionParameters([]byte{24, 43, 11}).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionId).
		Build(mockClient).
		Sign(privateKey)

	// note: yes this is the best way to add a ` to a raw string literal
	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x:\017\n\002\030\005\020\215\001\030\220N\"\003\030+\013"
sigMap: <
  sigPair: <
    ed25519: "-\326\n47\016[(\234\376\"\337\333\310\303uB\264\3305g{#R\005\212\354w\246xQ\016u\020G\346\311\344B\336\236/\311\351\021\036p/*N` + "`" + `j\224\035\307N\030\210\227\354\344>\213\005"
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
contractCall: <
  contractID: <
    contractNum: 5
  >
  gas: 141
  amount: 10000
  functionParameters: "\030+\013"
>
`

	assert.Equal(t, txString, tx.String())
}
