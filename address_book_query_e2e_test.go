//go:build all || testnets
// +build all testnets

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAddressBookQueryUpdateAll(t *testing.T) {
	client, err := ClientFromConfig([]byte(`{"network":"previewnet"}`))
	require.NoError(t, err)
	client.SetMirrorNetwork(previewnetMirror)

	previewnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		SetMaxAttempts(5).
		Execute(client)
	require.NoError(t, err)
	require.Greater(t, len(previewnet.NodeAddresses), 0)

	client, err = ClientFromConfig([]byte(`{"network":"testnet"}`))
	require.NoError(t, err)
	client.SetMirrorNetwork(testnetMirror)

	testnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		SetMaxAttempts(5).
		Execute(client)

	require.NoError(t, err)
	require.Greater(t, len(testnet.NodeAddresses), 0)

	client, err = ClientFromConfig([]byte(`{"network":"mainnet"}`))
	require.NoError(t, err)
	client.SetMirrorNetwork(mainnetMirror)

	mainnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		SetMaxAttempts(5).
		Execute(client)
	require.NoError(t, err)
	require.Greater(t, len(mainnet.NodeAddresses), 0)

	filePreviewnet, err := os.OpenFile("addressbook/previewnet.pb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	fileTestnet, err := os.OpenFile("addressbook/testnet.pb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	fileMainnet, err := os.OpenFile("addressbook/mainnet.pb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	_, err = filePreviewnet.Write(previewnet.ToBytes())
	require.NoError(t, err)

	_, err = fileTestnet.Write(testnet.ToBytes())
	require.NoError(t, err)

	_, err = fileMainnet.Write(mainnet.ToBytes())
	require.NoError(t, err)

	err = filePreviewnet.Close()
	require.NoError(t, err)

	err = fileTestnet.Close()
	require.NoError(t, err)

	err = fileMainnet.Close()
	require.NoError(t, err)
}
