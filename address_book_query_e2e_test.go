//go:build testnets
// +build testnets

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAddressBookQueryUpdateAll(t *testing.T) {
	client := ClientForPreviewnet()
	// There are some limitation on requests: unexpected HTTP status code received from server: 429 (Too Many Requests)
	time.Sleep(time.Second * 5)
	previewnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		Execute(client)
	require.NoError(t, err)
	require.Greater(t, len(previewnet.NodeAddresses), 0)

	client = ClientForTestnet()
	// There are some limitation on requests: unexpected HTTP status code received from server: 429 (Too Many Requests)
	// Testnet specifically has a more aggresive rate limit
	time.Sleep(time.Second * 20)
	testnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		Execute(client)
	require.NoError(t, err)
	require.Greater(t, len(testnet.NodeAddresses), 0)

	client = ClientForMainnet()
	// There are some limitation on requests: unexpected HTTP status code received from server: 429 (Too Many Requests)
	time.Sleep(time.Second * 5)
	mainnet, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
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
