//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
)

func TestUnitTopicMessageQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	topicID, err := TopicIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	topicInfo := NewTopicMessageQuery().
		SetTopicID(topicID)

	err = topicInfo.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTopicMessageQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	topicID, err := TopicIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	topicInfo := NewTopicMessageQuery().
		SetTopicID(topicID)

	err = topicInfo.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTopicMessageQueryGet(t *testing.T) {
	t.Parallel()

	topicID := TopicID{Topic: 7}

	balance := NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Now()).
		SetMaxAttempts(23).
		SetEndTime(time.Now()).
		SetLimit(32).
		SetCompletionHandler(func() {}).
		SetErrorHandler(func(stat status.Status) {}).
		SetRetryHandler(func(err error) bool { return false })

	balance.GetTopicID()
	balance.GetStartTime()
	balance.GetMaxAttempts()
	balance.GetEndTime()
	balance.GetLimit()
}

func TestUnitTopicMessageQueryNothingSet(t *testing.T) {
	t.Parallel()

	balance := NewTopicMessageQuery()

	balance.GetTopicID()
	balance.GetStartTime()
	balance.GetMaxAttempts()
	balance.GetEndTime()
	balance.GetLimit()
}
