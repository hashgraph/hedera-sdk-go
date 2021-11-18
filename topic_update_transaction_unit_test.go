//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTopicUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	topicID, err := TopicIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	topicUpdate := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAutoRenewAccountID(accountID)

	err = topicUpdate._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTopicUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)
	topicID, err := TopicIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	topicUpdate := NewTopicUpdateTransaction().
		SetTopicID(topicID).
		SetAutoRenewAccountID(accountID)

	err = topicUpdate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}
