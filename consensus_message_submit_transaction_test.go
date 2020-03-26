package hedera

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeConsensusMessageSubmitTransaction(t *testing.T) {
	date := time.Unix(1554158542, 0)

	testTopicID := ConsensusTopicID{Topic: 99}

	key, err := Ed25519PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	tx, err := NewConsensusMessageSubmitTransaction().
		SetTopicID(testTopicID).
		SetMessage([]byte("Hello Hashgraph")).
		SetTransactionValidDuration(24 * time.Hour).
		SetNodeAccountID(AccountID{Account: 3}).
		SetTransactionID(TransactionID{
			AccountID:  AccountID{Account: 2},
			ValidStart: date,
		}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		Build(nil)

	assert.NoError(t, err)

	tx.Sign(key)

	assert.Equal(t, `bodyBytes:"\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\030\300\204=\"\004\010\200\243\005\332\001\025\n\002\030c\022\017HelloHashgraph"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\307\333\307\010\\D`+"`"+`R\372\322\255zI$%\006\024\214\334\350\006g\021=\237\r\254e;+\234Y\335\235\246\270\215\234\235v\206e~F\261\025.\251yR\305s\301\347_\264\206\002XT\031<\004\002">>transactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:86400>consensusSubmitMessage:<topicID:<topicNum:99>message:"HelloHashgraph">`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}
