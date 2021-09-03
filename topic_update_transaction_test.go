package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const oldTopicMemo = "go-sdk::TestConsensusTopicUpdateTransaction_Execute::initial"

func TestIntegrationTopicUpdateTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(oldTopicMemo).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NoError(t, err)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, oldTopicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	newTopicMemo := "go-sdk::TestConsensusTopicUpdateTransaction_Execute::updated"

	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTopicMemo(newTopicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err = NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, newTopicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestUnitTopicUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	topicID, err := TopicIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	topicUpdate := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAutoRenewAccountID(accountID)

	err = topicUpdate._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTopicUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)
	topicID, err := TopicIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	topicUpdate := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAutoRenewAccountID(accountID)

	err = topicUpdate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch; some IDs have different networks set", err.Error())
	}
}

func TestIntegrationTopicUpdateTransactionNoMemo(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(oldTopicMemo).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NoError(t, err)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, oldTopicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	resp, err = NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTopicUpdateTransactionNoTopicID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(oldTopicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NoError(t, err)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, oldTopicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	_, err = NewTopicUpdateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	resp, err = NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
