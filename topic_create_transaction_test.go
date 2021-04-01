package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTopicCreateTransaction_Execute(t *testing.T) {
	client := newTestClient(t, false)

	topicMemo := "go-sdk::TestConsensusTopicCreateTransaction_Execute"

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetSubmitKey(client.GetOperatorPublicKey()).
		SetTopicMemo(topicMemo).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, client.GetOperatorPublicKey().String(), info.AdminKey.String())

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_TopicCreate_DifferentKeys(t *testing.T) {
	client := newTestClient(t, false)

	topicMemo := "go-sdk::TestConsensusTopicCreateTransaction_Execute"

	keys := make([]PrivateKey, 2)
	pubKeys := make([]PublicKey, 2)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	tx, err := NewTopicCreateTransaction().
		SetAdminKey(pubKeys[0]).
		SetSubmitKey(pubKeys[1]).
		SetTopicMemo(topicMemo).
		FreezeWith(client)
	assert.NoError(t, err)

	tx, err = tx.SignWithOperator(client)
	assert.NoError(t, err)
	tx.Sign(keys[0])
	resp, err := tx.Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, pubKeys[0].String(), info.AdminKey.String())

	txDelete, err := NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = txDelete.Sign(keys[0]).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_TopicCreate_JustSetMemo(t *testing.T) {
	client := newTestClient(t, false)

	topicMemo := "go-sdk::TestConsensusTopicCreateTransaction_Execute"

	resp, err := NewTopicCreateTransaction().
		SetTopicMemo(topicMemo).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
