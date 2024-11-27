//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

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

	mirrorNodeAddress, err := _ManagedNodeAddressFromString("hcs.mainnet.mirrornode.hedera.com:50211")
	require.NoError(t, err)
	require.True(t, *mirrorNodeAddress.address == "hcs.mainnet.mirrornode.hedera.com")
	require.True(t, mirrorNodeAddress.port == 50211)
	require.True(t, mirrorNodeAddress._String() == "hcs.mainnet.mirrornode.hedera.com:50211")

	mirrorNodeAddressSecure := mirrorNodeAddress._ToSecure()
	require.True(t, *mirrorNodeAddressSecure.address == "hcs.mainnet.mirrornode.hedera.com")
	require.True(t, mirrorNodeAddressSecure.port == 50212)
	require.True(t, mirrorNodeAddressSecure._String() == "hcs.mainnet.mirrornode.hedera.com:50212")

	mirrorNodeAddressInsecure := mirrorNodeAddressSecure._ToInsecure()
	require.True(t, *mirrorNodeAddressInsecure.address == "hcs.mainnet.mirrornode.hedera.com")
	require.True(t, mirrorNodeAddressInsecure.port == 50211)
	require.True(t, mirrorNodeAddressInsecure._String() == "hcs.mainnet.mirrornode.hedera.com:50211")

	_, err = _ManagedNodeAddressFromString("this is a random string with spaces:443")
	require.Error(t, err)

	_, err = _ManagedNodeAddressFromString("hcs.mainnet.mirrornode.hedera.com:notarealport")
	require.Error(t, err)
}
