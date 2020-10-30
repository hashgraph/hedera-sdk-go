package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestSerializeConsensusTopicDeleteTransaction(t *testing.T) {
// 	date := time.Unix(1554158542, 0)

// 	testTopicID := TopicID{Topic: 99}

// 	key, err := PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
// 	assert.NoError(t, err)

// 	tx, err := NewTopicDeleteTransaction().
// 		SetTopicID(testTopicID).
// 		SetTransactionValidDuration(24 * time.Hour).
// 		SetNodeAccountID(AccountID{Account: 3}).
// 		SetTransactionID(TransactionID{
// 			AccountID:  AccountID{Account: 2},
// 			ValidStart: date,
// 		}).
// 		SetMaxTransactionFee(HbarFromTinybar(1e6)).
// 		Freeze()

// 	assert.NoError(t, err)

// 	tx.Sign(key)

// 	assert.Equal(t, `bodyBytes:"\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\030\300\204=\"\004\010\200\243\005\322\001\004\n\002\030c"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\340\240!e\360j(\230Y\352q\213\334,\335`+"`"+`\236v\247\326T\026c-\373&\024\004~\036P\304\024u\211\220\370\306R\267\024<\372\3220\272\205?\335J^:=\232\021NU2\207\315r\024\317\013">>transactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:86400>consensusDeleteTopic:<topicID:<topicNum:99>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
// }

func TestConsensusTopicDeleteTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	topicMemo := "go-sdk::TestConsensusTopicDeleteTransaction_Execute"

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetTopicMemo(topicMemo).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	_, err = NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(22)).
		Execute(client)
	assert.NoError(t, err)

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	_, err = NewTopicInfoQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTopicID(topicID).
		SetQueryPayment(NewHbar(22)).
		Execute(client)
	assert.Error(t, err)

	status := err.(ErrHederaPreCheckStatus).Status
	assert.Equal(t, StatusInvalidTopicID, status)
}
