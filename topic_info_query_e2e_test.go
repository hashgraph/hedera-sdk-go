//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

const topicMemo = "go-sdk::topic memo"

func TestIntegrationTopicInfoQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	txID, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txID.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTopicInfoQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicInfo := NewTopicInfoQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTopicID(topicID)

	cost, err := topicInfo.GetCost(env.Client)
	require.NoError(t, err)

	_, err = topicInfo.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTopicInfoQuerySetBigMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicInfo := NewTopicInfoQuery().
		SetMaxQueryPayment(NewHbar(100000)).
		SetTopicID(topicID)

	cost, err := topicInfo.GetCost(env.Client)
	require.NoError(t, err)

	_, err = topicInfo.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTopicInfoQuerySetSmallMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicInfo := NewTopicInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetTopicID(topicID)

	cost, err := topicInfo.GetCost(env.Client)
	require.NoError(t, err)

	_, err = topicInfo.Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "cost of TopicInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tℏ", err.Error())
	}

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTopicInfoQueryInsufficientFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicInfo := NewTopicInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxQueryPayment(NewHbar(1)).
		SetTopicID(topicID)

	_, err = topicInfo.GetCost(env.Client)
	require.NoError(t, err)

	_, err = topicInfo.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INSUFFICIENT_TX_FEE", err.Error())
	}

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTopicInfoQueryThreshold(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	keys := make([]PrivateKey, 3)
	pubKeys := make([]PublicKey, 3)

	for i := range keys {
		newKey, err := PrivateKeyGenerateEd25519()
		if err != nil {
			panic(err)
		}

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	thresholdKey := KeyListWithThreshold(2).
		AddAllPublicKeys(pubKeys)

	txID, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetSubmitKey(thresholdKey).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := txID.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())
	assert.NotEmpty(t, info.SubmitKey.String())

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTopicInfoQueryNoTopicID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewTopicInfoQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INVALID_TOPIC_ID", err.Error())
	}

}
