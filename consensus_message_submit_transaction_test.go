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

    assert.Equal(t, `bodyBytes:"\n\014\n\006\010\316\247\212\345\005\022\002\030\002\022\002\030\003\030\300\204=\"\004\010\200\243\005\332\001)\n\002\030c\022\017HelloHashgraph\032\022\n\014\n\006\010\316\247\212\345\005\022\002\030\002\020\001\030\001"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"|T\235GE\003\210\027\302\231]\002\265\304\376\034E\224\254\376\225\266\006|b\025s\357W\236_9\217^\234^\230\235el8\032H$\355y\232\327J)\213\236\3669\016\235\300\342\030\323'\366\201\n">>transactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:86400>consensusSubmitMessage:<topicID:<topicNum:99>message:"HelloHashgraph"chunkInfo:<initialTransactionID:<transactionValidStart:<seconds:1554158542>accountID:<accountNum:2>>total:1number:1>>`, strings.ReplaceAll(strings.ReplaceAll(tx.List[0].String(), " ", ""), "\n", ""))
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
