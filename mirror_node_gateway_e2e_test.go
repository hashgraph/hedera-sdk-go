//go:build all || e2e
// +build all e2e

package hedera

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
const mainnetMirrorNodeUrl = "mainnet-public.mirrornode.hedera.com"
const testnetMirrorNodeUrl = "testnet.mirrornode.hedera.com"
const previewnetMirrorNodeUrl = "previewnet.mirrornode.hedera.com"

func TestAccountInfoTestnet(t *testing.T) {
	testAccountInfoQuery(t, testnetMirrorNodeUrl)
}
func TestAccountInfoMainnet(t *testing.T) {
	testAccountInfoQuery(t, mainnetMirrorNodeUrl)
}
func TestAccountInfoPreviewnet(t *testing.T) {
	testAccountInfoQuery(t, previewnetMirrorNodeUrl)
}
func testAccountInfoQuery(t *testing.T, network string) {
	t.Parallel()

	result, e := accountInfoQuery(network, "1")
	require.NoError(t, e)
	assert.Equal(t, 20, len(result))
}

func TestAccountBalanceTestnet(t *testing.T) {
	testAccountBalanceQuery(t, testnetMirrorNodeUrl)
}
func TestAccountBalanceMainnet(t *testing.T) {
	testAccountBalanceQuery(t, mainnetMirrorNodeUrl)
}
func TestAccountBalancePreviewnet(t *testing.T) {
	testAccountBalanceQuery(t, previewnetMirrorNodeUrl)
}
func testAccountBalanceQuery(t *testing.T, network string) {
	t.Parallel()

	result, e := accountBalanceQuery(network, "1")
	require.NoError(t, e)
	assert.Equal(t, 3, len(result))
	_, exist := result["balance"]
	require.True(t, exist)
	_, exist = result["timestamp"]
	require.True(t, exist)
	_, exist = result["tokens"]
	require.True(t, exist)
}

func TestContractInfoPreviewnetContractNotFound(t *testing.T) {
	t.Parallel()

	result, e := contractInfoQuery(previewnetMirrorNodeUrl, "1")
	require.Error(t, e)
	assert.True(t, result == nil)
}
func TestContractInfoTestnet(t *testing.T) {
	t.Parallel()

	result, e := contractInfoQuery(testnetMirrorNodeUrl, "0.0.7376843")
	require.NoError(t, e)
	_, exist := result["bytecode"]
	require.True(t, exist)
}

func TestBuildUrlReturnCorrectUrl(t *testing.T) {
	url := "https://testnet.mirrornode.hedera.com/api/v1/accounts/0.0.7477022/tokens"

	result := buildUrl("testnet.mirrornode.hedera.com", "accounts", "0.0.7477022", "tokens")
	assert.Equal(t, url, result)
}
