package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeSystemDeleteFileIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewSystemDeleteTransaction().
		SetFileID(FileID{File: 3}).
		SetExpirationTime(time.Unix(15415151511, 0)).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		SetTransactionID(testTransactionID).
		FreezeWith(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\240\215\006\"\002\010x\242\001\014\032\006\010\227\227\302\2669\n\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\237\271g\353\270aW\256TM\t;\362P\312v\220\261\003P\272\032\326x_2\345\371\233RgV6\362\252\341\357/\177\271m\373\223\233\000\204\326\347m|\232\344\334%M\325Mh\374\035\351\345\016\013">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000transactionValidDuration:<seconds:120>systemDelete:<fileID:<fileNum:3>expirationTime:<seconds:15415151511>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestSerializeSystemDeleteContractIDTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewSystemDeleteTransaction().
		SetContractID(ContractID{Contract: 3}).
		SetExpirationTime(time.Unix(15415151511, 0)).
		SetMaxTransactionFee(HbarFromTinybar(100_000)).
		SetTransactionID(testTransactionID).
		FreezeWith(mockClient)

	assert.NoError(t, err)
	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\240\215\006\"\002\010x\242\001\014\032\006\010\227\227\302\2669\022\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"F\317\270L\206\376\2764\247\326\3104\300\227\207\375\231T\300\203AVAu)\000\347\\\013@\250\246\374pU\306\311O\212\354}1\355\3328\272\333\206\023D\322\023\354\"\235\232:\024\033\365\203\207i\005">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:100000transactionValidDuration:<seconds:120>systemDelete:<contractID:<contractNum:3>expirationTime:<seconds:15415151511>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}
