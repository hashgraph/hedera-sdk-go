//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitNetworkAddressBookGetsSet(t *testing.T) {
	network := _NewNetwork()
	network._SetTransportSecurity(true)

	ledgerID, err := LedgerIDFromString("mainnet")
	require.NoError(t, err)

	network._SetLedgerID(*ledgerID)
	require.NoError(t, err)

	require.True(t, network.addressBook != nil)
}
