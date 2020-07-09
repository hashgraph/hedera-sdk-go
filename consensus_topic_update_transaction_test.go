package hedera

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeConsensusTopicUpdateTransaction(t *testing.T) {
	date := time.Unix(1554158542, 0)

	testTopicID := ConsensusTopicID{Topic: 99}

	key, err := Ed25519PrivateKeyFromString("302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962")
	assert.NoError(t, err)

	tx, err := NewConsensusTopicUpdateTransaction().
		SetTopicID(testTopicID).SetAdminKey(key.PublicKey()).SetTopicMemo("updated topic memo").
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

	assert.Equal(t, `bodyBytes:"\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\030\300\204=\"\004\010\200\243\005\312\001>\n\002\030c\022\024\n\022updatedtopicmemo2\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\331\274\230\260l\"l\322\265\245\225\265\r\221\250N5\316gSL\004\003\237X\325\220\226\314O$\344poH\322\330\226\013\021\253.\222i\355\202\303\360\016\3036\213'\003T\033\2121\213\376\2650\010">>transactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:86400>consensusUpdateTopic:<topicID:<topicNum:99>memo:<value:"updatedtopicmemo">adminKey:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestConsensusTopicUpdateTransaction_Execute(t *testing.T) {
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

	oldTopicMemo := "go-sdk::TestConsensusTopicUpdateTransaction_Execute::initial"

	txID, err := NewConsensusTopicCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetTopicMemo(oldTopicMemo).
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

	assert.Equal(t, oldTopicMemo, info.Memo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, client.GetOperatorKey().String(), info.AdminKey.String())

	newTopicMemo := "go-sdk::TestConsensusTopicUpdateTransaction_Execute::updated"

	txID, err = NewConsensusTopicUpdateTransaction().
		SetTopicID(topicID).
		SetTopicMemo(newTopicMemo).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	info, err = NewConsensusTopicInfoQuery().
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, newTopicMemo, info.Memo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, client.GetOperatorKey().String(), info.AdminKey.String())

	_, err = NewConsensusTopicDeleteTransaction().
		SetTopicID(topicID).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
}
