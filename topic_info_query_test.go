package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const topicMemo = "go-sdk::topic memo"

func TestIntegrationTopicInfoQueryCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	txID, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestUnitTopicInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	topicID, err := TopicIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	topicInfo := NewTopicInfoQuery().
		SetTopicID(topicID)

	err = topicInfo._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTopicInfoQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	topicID, err := TopicIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	topicInfo := NewTopicInfoQuery().
		SetTopicID(topicID)

	err = topicInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestIntegrationTopicInfoQueryGetCost(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicInfo := NewTopicInfoQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetTopicID(topicID)

	cost, err := topicInfo.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = topicInfo.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTopicInfoQuerySetBigMaxPayment(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicInfo := NewTopicInfoQuery().
		SetMaxQueryPayment(NewHbar(100000)).
		SetTopicID(topicID)

	cost, err := topicInfo.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = topicInfo.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTopicInfoQuerySetSmallMaxPayment(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicInfo := NewTopicInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetTopicID(topicID)

	cost, err := topicInfo.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = topicInfo.Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "cost of TopicInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tℏ", err.Error())
	}

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTopicInfoQueryInsufficientFee(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicInfo := NewTopicInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetMaxQueryPayment(NewHbar(1)).
		SetTopicID(topicID)

	_, err = topicInfo.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = topicInfo.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INSUFFICIENT_TX_FEE", err.Error())
	}

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTopicInfoQueryThreshold(t *testing.T) {
	env := NewIntegrationTestEnv(t)

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

	txID, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetSubmitKey(thresholdKey).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		SetMaxQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, env.Client.GetOperatorPublicKey().String(), info.AdminKey.String())
	assert.NotEmpty(t, info.SubmitKey.String())

	_, err = NewTopicDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicID(topicID).
		Execute(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestIntegrationTopicInfoQueryNoTopicID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewTopicInfoQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INVALID_TOPIC_ID", err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
