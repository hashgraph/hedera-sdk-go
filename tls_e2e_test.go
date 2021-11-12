//+build all e2e

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationPreviewnetTls(t *testing.T) {
	var network = map[string]AccountID{
		"0.previewnet.hedera.com:50212": {Account: 3},
		"1.previewnet.hedera.com:50212": {Account: 4},
		"2.previewnet.hedera.com:50212": {Account: 5},
		"3.previewnet.hedera.com:50212": {Account: 6},
		"4.previewnet.hedera.com:50212": {Account: 7},
	}

	client := ClientForNetwork(network)
	client.SetNetworkName(NetworkNamePreviewnet)
	client.SetMirrorNetwork([]string{"hcs.previewnet.mirrornode.hedera.com:5600"})
	client.SetMaxAttempts(3)

	for _, nodeAccountID := range network {
		_, err := NewAccountBalanceQuery().
			SetNodeAccountIDs([]AccountID{nodeAccountID}).
			SetAccountID(nodeAccountID).
			Execute(client)
		assert.NoError(t, err)
	}
}

func TestIntegrationTestnetTls(t *testing.T) {
	var network = map[string]AccountID{
		"0.testnet.hedera.com:50212": {Account: 3},
		"1.testnet.hedera.com:50212": {Account: 4},
		"2.testnet.hedera.com:50212": {Account: 5},
		"3.testnet.hedera.com:50212": {Account: 6},
		"4.testnet.hedera.com:50212": {Account: 7},
	}

	client := ClientForNetwork(network)
	client.SetNetworkName(NetworkNameTestnet)
	client.SetMirrorNetwork([]string{"hcs.testnet.mirrornode.hedera.com:5600"})
	client.SetMaxAttempts(3)

	for _, nodeAccountID := range network {
		_, err := NewAccountBalanceQuery().
			SetNodeAccountIDs([]AccountID{nodeAccountID}).
			SetAccountID(nodeAccountID).
			Execute(client)
		assert.NoError(t, err)
	}
}
