//go:build all || unit
// +build all unit

package hedera

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

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitManagedNodeAddressTest(t *testing.T) {
	t.Parallel()

	ipAddress, err := _ManagedNodeAddressFromString("35.237.200.180:50211")
	require.NoError(t, err)
	require.True(t, *ipAddress.address == "35.237.200.180")
	require.True(t, ipAddress.port == 50211)
	require.True(t, ipAddress._String() == "35.237.200.180:50211")

	ipAddressSecure := ipAddress._ToSecure()
	require.True(t, *ipAddressSecure.address == "35.237.200.180")
	require.True(t, ipAddressSecure.port == 50212)
	require.True(t, ipAddressSecure._String() == "35.237.200.180:50212")

	ipAddressInsecure := ipAddressSecure._ToInsecure()
	require.True(t, *ipAddressInsecure.address == "35.237.200.180")
	require.True(t, ipAddressInsecure.port == 50211)
	require.True(t, ipAddressInsecure._String() == "35.237.200.180:50211")

	urlAddress, err := _ManagedNodeAddressFromString("0.testnet.hedera.com:50211")
	require.NoError(t, err)
	require.True(t, *urlAddress.address == "0.testnet.hedera.com")
	require.True(t, urlAddress.port == 50211)
	require.True(t, urlAddress._String() == "0.testnet.hedera.com:50211")

	urlAddressSecure := urlAddress._ToSecure()
	require.True(t, *urlAddressSecure.address == "0.testnet.hedera.com")
	require.True(t, urlAddressSecure.port == 50212)
	require.True(t, urlAddressSecure._String() == "0.testnet.hedera.com:50212")

	urlAddressInsecure := urlAddressSecure._ToInsecure()
	require.True(t, *urlAddressInsecure.address == "0.testnet.hedera.com")
	require.True(t, urlAddressInsecure.port == 50211)
	require.True(t, urlAddressInsecure._String() == "0.testnet.hedera.com:50211")

	mirrorNodeAddress, err := _ManagedNodeAddressFromString("hcs.mainnet.mirrornode.hedera.com:5600")
	require.NoError(t, err)
	require.True(t, *mirrorNodeAddress.address == "hcs.mainnet.mirrornode.hedera.com")
	require.True(t, mirrorNodeAddress.port == 5600)
	require.True(t, mirrorNodeAddress._String() == "hcs.mainnet.mirrornode.hedera.com:5600")

	mirrorNodeAddressSecure := mirrorNodeAddress._ToSecure()
	require.True(t, *mirrorNodeAddressSecure.address == "hcs.mainnet.mirrornode.hedera.com")
	require.True(t, mirrorNodeAddressSecure.port == 443)
	require.True(t, mirrorNodeAddressSecure._String() == "hcs.mainnet.mirrornode.hedera.com:443")

	mirrorNodeAddressInsecure := mirrorNodeAddressSecure._ToInsecure()
	require.True(t, *mirrorNodeAddressInsecure.address == "hcs.mainnet.mirrornode.hedera.com")
	require.True(t, mirrorNodeAddressInsecure.port == 5600)
	require.True(t, mirrorNodeAddressInsecure._String() == "hcs.mainnet.mirrornode.hedera.com:5600")

	_, err = _ManagedNodeAddressFromString("this is a random string with spaces:443")
	require.Error(t, err)

	_, err = _ManagedNodeAddressFromString("hcs.mainnet.mirrornode.hedera.com:notarealport")
	require.Error(t, err)
}
