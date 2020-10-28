package hedera

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestSerializeTopicCreateTransaction(t *testing.T) {
// 	date := time.Unix(1554158542, 0)

// 	key, err := PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
// 	assert.NoError(t, err)

// 	tx, err := NewTopicCreateTransaction().
// 		SetTopicMemo("this is a test topic").
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

// 	assert.Equal(t, `bodyBytes:"\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\030\300\204=\"\004\010\200\243\005\302\001\035\n\024thisisatesttopic2\005\010\320\310\341\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"Dqt!\253\207\035\227\224\032/WTV\023\313H\357Z\220\357^[\270\325\361\340).\233\326(>\324\303\332\244\033Z>\240\206\301\357\213Wr\033\321\34095\252\\`+"`"+`\254qa\"\007\014">>transactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:86400>consensusCreateTopic:<memo:"thisisatesttopic"autoRenewPeriod:<seconds:7890000>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
// }

func TestTopicCreateTransaction_Execute(t *testing.T) {
	var client *Client

	network := os.Getenv("HEDERA_NETWORK")

	if network == "previewnet" {
		client = ClientForPreviewnet()
	}

	client, err := ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	topicMemo := "go-sdk::TestConsensusTopicCreateTransaction_Execute"

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

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs(nodeIDs).
		SetQueryPayment(NewHbar(22)).
		Execute(client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.Memo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, client.GetOperatorKey().String(), info.AdminKey.String())

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs(nodeIDs).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
