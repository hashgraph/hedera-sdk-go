package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeTopicInfoQuery(t *testing.T) {
	query := NewTopicInfoQuery().
		SetTopicID(TopicID{Topic: 3}).
		Query

	assert.Equal(t, `consensusGetTopicInfo:{header:{}topicID:{topicNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestTopicInfoQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	topicMemo := "go-sdk::TestConsensusTopicInfoQuery_Execute"

	txID, err := NewTopicCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetTopicMemo(topicMemo).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.Memo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, client.GetOperatorPublicKey().String(), info.AdminKey.String())

	_, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		Execute(client)
	assert.NoError(t, err)
}

func TestTopicInfoQuery_Threshold_Execute(t *testing.T) {
	client := newTestClient(t)

	keys := make([]PrivateKey, 3)
	pubKeys := make([]PublicKey, 3)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		if err != nil {
			panic(err)
		}

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	thresholdKey := KeyListWithThreshold(2).
		AddAllPublicKeys(pubKeys)

	topicMemo := "go-sdk::TestConsensusTopicInfoQuery_Execute"

	txID, err := NewTopicCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetSubmitKey(thresholdKey).
		SetTopicMemo(topicMemo).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.Memo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, client.GetOperatorPublicKey().String(), info.AdminKey.String())
	assert.NotEmpty(t, info.SubmitKey.String())

	_, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		Execute(client)
	assert.NoError(t, err)
}

func Test_TopicInfo_NoTopicID(t *testing.T) {
	client := newTestClient(t)

	_, err := NewTopicInfoQuery().
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_TOPIC_ID"), err.Error())
	}
}
