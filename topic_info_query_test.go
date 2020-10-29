package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsensusTopicInfoQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	topicMemo := "go-sdk::TestConsensusTopicInfoQuery_Execute"

	txID, err := NewTopicCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetTopicMemo(topicMemo).
		SetMaxTransactionFee(NewHbar(1)).
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
	assert.Equal(t, client.GetOperatorKey().String(), info.AdminKey.String())

	_, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetMaxTransactionFee(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)
}
