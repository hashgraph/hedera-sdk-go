//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
