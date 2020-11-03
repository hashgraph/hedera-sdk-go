package hedera

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTopicMessageQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	wait := true
	start := time.Now()

	_, err = NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(client, func(message TopicMessage) {
			if string(message.Contents) == "Hello from Hedera SDK" {
				wait = false
			}
		})
	assert.NoError(t, err)

	resp, err = NewTopicMessageSubmitTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMessage([]byte("Hello from Hedera SDK")).
		SetTopicID(topicID).
		Execute(client)
	assert.NoError(t, err)

	for {
		if !wait || uint64(time.Since(start).Seconds()) > 30 {
			break
		}

		time.Sleep(2500)
	}

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	if wait {
		panic("Message was not received within 30 seconds")
	}
}
