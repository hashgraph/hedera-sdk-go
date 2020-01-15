package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeSystemUndeleteFileIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewSystemUndeleteTransaction().
		SetFileID(FileID{File: 3}).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\252\001\004\n\002\030\003"
sigMap: <
  sigPair: <
    pubKeyPrefix: "\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"
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
 `

	assert.Equal(t, txString, tx.String())
}

func TestSerializeSystemUndeleteContractIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx := NewSystemUndeleteTransaction().
		SetContractID(ContractID{Contract: 3}).
		SetMaxTransactionFee(1e6).
		SetTransactionID(testTransactionID).
		Build(mockClient).
		Sign(privateKey)

	txString := `bodyBytes: "\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\252\001\004\022\002\030\003"
sigMap: <
  sigPair: <
    pubKeyPrefix: "\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"
    ed25519: "\271\205<9'>\341\013\245\026f\005^\252\210*\262\315\030>;\330-\351\212 \313\000\265\037\315 \243\363\206\256\316\227\233\203\225\206\333\363ik^\310\264p\377\306\272\247j\361\026\246\302\264|\025\203\013"
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
  contractID: <
    contractNum: 3
  >
>
`

	assert.Equal(t, txString, tx.String())
}
