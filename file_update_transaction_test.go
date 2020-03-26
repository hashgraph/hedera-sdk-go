package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeFileUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewFileUpdateTransaction().
		SetFileID(FileID{File: 5}).
		SetContents([]byte("there was a hole here")).
		SetExpirationTime(time.Unix(15415151511, 0)).
		AddKey(privateKey.PublicKey()).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)
	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\232\001I\n\002\030\005\022\006\010\227\227\302\2669\032$\n\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216\"\025therewasaholehere"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\243.\202\276AZ\273i4\342\262a+:\231a\350;\326%\016\304\314\271b\205\261\316l\214bot\304+5\241\034N\r\361\340\031\360OZ\356\0149\2321\037\377\232\3515\324o\303\316\243\237\017">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>fileUpdate:<fileID:<fileNum:5>expirationTime:<seconds:15415151511>keys:<keys:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">>contents:"therewasaholehere">`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}
