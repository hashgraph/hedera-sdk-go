//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func DisabledTestIntegrationFeeSchedulesFromBytes(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)

	feeSchedulesBytes, err := NewFileContentsQuery().
		SetFileID(FileID{Shard: 0, Realm: 0, File: 111}).
		Execute(env.Client)
	require.NoError(t, err)
	feeSchedules, err := FeeSchedulesFromBytes(feeSchedulesBytes)
	require.NoError(t, err)
	assert.NotNil(t, feeSchedules)
	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].FeeData.NodeData.Constant, int64(4498129603))
	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].FeeData.ServiceData.Constant, int64(71970073651))
	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].RequestType, RequestTypeCryptoCreate)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func DisabledTestIntegrationNodeAddressBookFromBytes(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)

	nodeAddressBookBytes, err := NewFileContentsQuery().
		SetFileID(FileID{Shard: 0, Realm: 0, File: 101}).
		Execute(env.Client)
	require.NoError(t, err)
	nodeAddressbook, err := _NodeAddressBookFromBytes(nodeAddressBookBytes)
	require.NoError(t, err)
	assert.NotNil(t, nodeAddressbook)

	for _, ad := range nodeAddressbook.nodeAddresses {
		println(ad.nodeID)
		println(string(ad.certHash))
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
