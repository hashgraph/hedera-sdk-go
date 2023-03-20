//go:build all || unit
// +build all unit

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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitClientFromConfig(t *testing.T) {
	client, err := ClientFromConfig([]byte(testClientJSON))
	require.NoError(t, err)

	assert.NotNil(t, client)
	assert.True(t, len(client.network.network) > 0)
	assert.Nil(t, client.operator)
}

func TestUnitClientFromConfigWithOperator(t *testing.T) {
	client, err := ClientFromConfig([]byte(testClientJSONWithOperator))
	require.NoError(t, err)

	assert.NotNil(t, client)

	testOperatorKey, err := PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	require.NoError(t, err)

	assert.True(t, len(client.network.network) > 0)
	assert.NotNil(t, client.operator)
	assert.Equal(t, testOperatorKey.ed25519PrivateKey.keyData, client.operator.privateKey.ed25519PrivateKey.keyData)
	assert.Equal(t, AccountID{Account: 3}.Account, client.operator.accountID.Account)
}

func TestUnitClientFromConfigWrongType(t *testing.T) {
	_, err := ClientFromConfig([]byte(testClientJSONWrongTypeMirror))
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "mirrorNetwork is expected to be either string or an array of strings", err.Error())
	}

	_, err = ClientFromConfig([]byte(testClientJSONWrongTypeNetwork))
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network is expected to be map of string to string, or string", err.Error())
	}
}

func TestUnitClientSetNetworkExtensive(t *testing.T) {
	client := ClientForTestnet()
	nodes := make(map[string]AccountID, 2)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}

	err := client.SetNetwork(nodes)
	require.NoError(t, err)
	network := client.GetNetwork()
	assert.Equal(t, 2, len(network))
	assert.Equal(t, network["0.testnet.hedera.com:50211"], AccountID{0, 0, 3, nil, nil, nil})
	assert.Equal(t, network["1.testnet.hedera.com:50211"], AccountID{0, 0, 4, nil, nil, nil})

	nodes = make(map[string]AccountID, 2)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["2.testnet.hedera.com:50211"] = AccountID{0, 0, 5, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	network = client.GetNetwork()
	assert.Equal(t, 3, len(network))
	assert.Equal(t, network["0.testnet.hedera.com:50211"], AccountID{0, 0, 3, nil, nil, nil})
	assert.Equal(t, network["1.testnet.hedera.com:50211"], AccountID{0, 0, 4, nil, nil, nil})
	assert.Equal(t, network["2.testnet.hedera.com:50211"], AccountID{0, 0, 5, nil, nil, nil})

	nodes = make(map[string]AccountID, 1)
	nodes["2.testnet.hedera.com:50211"] = AccountID{0, 0, 5, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	network = client.GetNetwork()
	networkMirror := client.GetMirrorNetwork()
	assert.Equal(t, 1, len(network))
	assert.Equal(t, network["2.testnet.hedera.com:50211"], AccountID{0, 0, 5, nil, nil, nil})
	// There is only one mirror address, no matter the transport security
	assert.Equal(t, networkMirror[0], "testnet.mirrornode.hedera.com:443")

	client.SetTransportSecurity(true)
	client.SetCertificateVerification(true)
	network = client.GetNetwork()
	networkMirror = client.GetMirrorNetwork()
	assert.Equal(t, network["2.testnet.hedera.com:50212"], AccountID{0, 0, 5, nil, nil, nil})
	assert.Equal(t, networkMirror[0], "testnet.mirrornode.hedera.com:443")

	err = client.Close()
	require.NoError(t, err)
}

func TestUnitClientSetMirrorNetwork(t *testing.T) {
	defaultNetwork := make([]string, 0)
	defaultNetwork = append(defaultNetwork, "hcs.testnet.mirrornode.hedera.com:5600")
	client := ClientForTestnet()
	client.SetMirrorNetwork(defaultNetwork)

	mirrorNetwork := client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, "hcs.testnet.mirrornode.hedera.com:5600", mirrorNetwork[0])

	defaultNetworkWithExtraNode := make([]string, 0)
	defaultNetworkWithExtraNode = append(defaultNetworkWithExtraNode, "hcs.testnet.mirrornode.hedera.com:5600")
	defaultNetworkWithExtraNode = append(defaultNetworkWithExtraNode, "hcs.testnet1.mirrornode.hedera.com:5600")

	client.SetMirrorNetwork(defaultNetworkWithExtraNode)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 2, len(mirrorNetwork))
	require.True(t, contains(mirrorNetwork, "hcs.testnet.mirrornode.hedera.com:5600"))
	require.True(t, contains(mirrorNetwork, "hcs.testnet1.mirrornode.hedera.com:5600"))

	defaultNetwork = make([]string, 0)
	defaultNetwork = append(defaultNetwork, "hcs.testnet1.mirrornode.hedera.com:5600")

	client.SetMirrorNetwork(defaultNetwork)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, "hcs.testnet1.mirrornode.hedera.com:5600", mirrorNetwork[0])

	defaultNetwork = make([]string, 0)
	defaultNetwork = append(defaultNetwork, "hcs.testnet.mirrornode.hedera.com:5600")

	client.SetMirrorNetwork(defaultNetwork)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, "hcs.testnet.mirrornode.hedera.com:5600", mirrorNetwork[0])

	client.SetTransportSecurity(true)
	mirrorNetwork = client.GetMirrorNetwork()
	// SetTransportSecurity is deprecated, so the mirror node should not be updated
	assert.Equal(t, "hcs.testnet.mirrornode.hedera.com:5600", mirrorNetwork[0])

	err := client.Close()
	require.NoError(t, err)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestUnitClientSetMultipleNetwork(t *testing.T) {
	client := ClientForTestnet()
	nodes := make(map[string]AccountID, 8)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["34.94.106.61:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["50.18.132.211:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["138.91.142.219:50211"] = AccountID{0, 0, 3, nil, nil, nil}

	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["35.237.119.55:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["3.212.6.13:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["52.168.76.241:50211"] = AccountID{0, 0, 4, nil, nil, nil}

	err := client.SetNetwork(nodes)
	require.NoError(t, err)
	net := client.GetNetwork()

	if val, ok := net["0.testnet.hedera.com:50211"]; ok {
		require.Equal(t, val.String(), "0.0.3")
	}

	if val, ok := net["1.testnet.hedera.com:50211"]; ok {
		require.Equal(t, val.String(), "0.0.4")
	}

	if val, ok := net["50.18.132.211:50211"]; ok {
		require.Equal(t, val.String(), "0.0.3")
	}

	if val, ok := net["3.212.6.13:50211"]; ok {
		require.Equal(t, val.String(), "0.0.4")
	}

}
