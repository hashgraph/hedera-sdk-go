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
