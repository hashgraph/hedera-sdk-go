package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationTopicCreateTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetSubmitKey(env.Client.GetOperatorPublicKey()).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
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

func TestIntegrationTopicCreateTransactionDifferentKeys(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 2)
	pubKeys := make([]PublicKey, 2)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	tx, err := NewTopicCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAdminKey(pubKeys[0]).
		SetSubmitKey(pubKeys[1]).
		SetTopicMemo(topicMemo).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	tx, err = tx.SignWithOperator(env.Client)
	assert.NoError(t, err)
	tx.Sign(keys[0])
	resp, err := tx.Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	info, err := NewTopicInfoQuery().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, topicMemo, info.TopicMemo)
	assert.Equal(t, uint64(0), info.SequenceNumber)
	assert.Equal(t, pubKeys[0].String(), info.AdminKey.String())

	txDelete, err := NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = txDelete.Sign(keys[0]).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestUnitTopicCreateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	topicCreate := NewTopicCreateTransaction().
		SetAutoRenewAccountID(accountID)

	err = topicCreate._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTopicCreateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	topicCreate := NewTopicCreateTransaction().
		SetAutoRenewAccountID(accountID)

	err = topicCreate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestIntegrationTopicCreateTransactionJustSetMemo(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTopicCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTopicMemo(topicMemo).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
