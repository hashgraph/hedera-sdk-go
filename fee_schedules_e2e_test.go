//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func DisabledTestIntegrationNodeAddressBookFromBytes(t *testing.T) { // nolint
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	nodeAddressBookBytes, err := NewFileContentsQuery().
		SetFileID(FileID{Shard: 0, Realm: 0, File: 101}).
		Execute(env.Client)
	require.NoError(t, err)
	nodeAddressbook, err := NodeAddressBookFromBytes(nodeAddressBookBytes)
	require.NoError(t, err)
	assert.NotNil(t, nodeAddressbook)

	for _, ad := range nodeAddressbook.NodeAddresses {
		println(ad.NodeID)
		println(string(ad.CertHash))
	}

}
