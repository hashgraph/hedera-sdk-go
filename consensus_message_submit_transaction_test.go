package hedera

import (
	"os"
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

    assert.Equal(t, `bodyBytes:"\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\"\002\010x\332\001)\n\002\030c\022\017HelloHashgraph\032\022\n\014\n\006\010\316\247\212\345\005\022\002\030\002\020\001\030\001"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\366No\306\344\222\020m\341\376}W\306\2257x\020B\024\256X\000\240s\257\314e\205\337t\325\256)&|\247\006I\016\032\210K\216\3273\246\265\036V\r\020\347\343f\207\346\0211\377\345\023\307J\016">>transactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>nodeAccountID:<accountNum:3>transactionValidDuration:<seconds:120>consensusSubmitMessage:<topicID:<topicNum:99>message:"HelloHashgraph"chunkInfo:<initialTransactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>total:1number:1>>`, strings.ReplaceAll(strings.ReplaceAll(tx.List[0].String(), " ", ""), "\n", ""))
}

func TestConsensusMessageSubmitTransaction_Execute(t *testing.T) {
	client, err := ClientFromFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := Ed25519PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	txID, err := NewConsensusTopicCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetTopicMemo("go-sdk::TestConsensusMessageSubmitTransaction_Execute").
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	topicID := receipt.GetConsensusTopicID()
	assert.NotNil(t, topicID)

	info, err := NewConsensusTopicInfoQuery().
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, uint64(0), info.SequenceNumber)

    txIDs, err := NewConsensusMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage([]byte("go-sdk::TestConsensusMessageSubmitTransaction_Execute::MessageSubmit")).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = txIDs[0].GetReceipt(client)
	assert.NoError(t, err)

	info, err = NewConsensusTopicInfoQuery().
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, uint64(1), info.SequenceNumber)

	txID, err = NewConsensusTopicDeleteTransaction().
		SetTopicID(topicID).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
