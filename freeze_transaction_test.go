package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"time"

	"testing"
)

func TestSerializeFreezeTransaction(t *testing.T) {
	startTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		12, 30, 0, time.Now().Nanosecond(), time.Now().Location(),
	)

	endTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		14, 30, 0, time.Now().Nanosecond(), time.Now().Location(),
	)

	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewFreezeTransaction().
		SetTransactionID(testTransactionID).
		SetStartTime(startTime).
		SetEndTime(endTime).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		FreezeWith(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\272\001\004\010\016\020\036"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"P\220\234@\324\370\027\217&\331\263\377\"9\337\302\013\315-jg\253uz\177\362\211~\331^:A\357\315S\317\320\024\275\343wH\246|\325\332$_J\376,\317<<\314_\361\364\317\3218\317<\006">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>freeze:<startHour:14startMin:30>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}
