package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"

	"testing"
)

func TestSerializeSystemUndeleteFileIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewSystemUndeleteTransaction().
		SetFileID(FileID{File: 3}).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		FreezeWith(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\240\215\006\"\002\010x\252\001\004\n\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\203\300S\261\351,\316\351-0\321s\233(k\301\027&0M\274\253O\301\022\214a\362\001L\365\333|\275\335/J\365\207A\373\177j\222\001TR\017\004(\016\213\236G\034\013\034l\321\273\201\273\010">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000transactionValidDuration:<seconds:120>systemUndelete:<fileID:<fileNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestSerializeSystemUndeleteContractIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewSystemUndeleteTransaction().
		SetContractID(ContractID{Contract: 3}).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		FreezeWith(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\240\215\006\"\002\010x\252\001\004\022\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\255\222V\361)VL\237\270{\274q\305\216\337\316\265z\216\320\344m\203\304K\357\205]\355\237\001P\026\324+z\237\353\332\3170F\206\3355\271Fm\216*\305\030\343\361\274x\327*\264\377\310\254\336\004">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000transactionValidDuration:<seconds:120>systemUndelete:<contractID:<contractNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}
