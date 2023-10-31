//go:build all || e2e
// +build all e2e

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

// func TestIntegrationPreviewnetTls(t *testing.T) {
// 	var network = map[string]AccountID{
// 		"0.previewnet.hedera.com:50212": {Account: 3},
// 		"1.previewnet.hedera.com:50212": {Account: 4},
// 		"2.previewnet.hedera.com:50212": {Account: 5},
// 		// "3.previewnet.hedera.com:50212": {Account: 6},
// 		"4.previewnet.hedera.com:50212": {Account: 7},
// 	}

// 	client := ClientForNetwork(network)
// 	ledger, _ := LedgerIDFromNetworkName(NetworkNamePreviewnet)
// 	client.SetTransportSecurity(true)
// 	client.SetLedgerID(*ledger)
// 	client.SetMaxAttempts(3)

// 	for _, nodeAccountID := range network {
// 		_, err := NewAccountBalanceQuery().
// 			SetNodeAccountIDs([]AccountID{nodeAccountID}).
// 			SetAccountID(nodeAccountID).
// 			Execute(client)
// 		require.NoError(t, err)
// 	}
// }

func TestIntegrationTestnetTls(t *testing.T) {
	var network = map[string]AccountID{
		"0.testnet.hedera.com:50212": {Account: 3},
		"1.testnet.hedera.com:50212": {Account: 4},
		"2.testnet.hedera.com:50212": {Account: 5},
		"3.testnet.hedera.com:50212": {Account: 6},
		"4.testnet.hedera.com:50212": {Account: 7},
	}

	client := ClientForNetwork(network)
	ledger, _ := LedgerIDFromNetworkName(NetworkNameTestnet)
	client.SetTransportSecurity(true)
	client.SetLedgerID(*ledger)
	client.SetMaxAttempts(3)

	for _, nodeAccountID := range network {
		_, err := NewAccountBalanceQuery().
			SetNodeAccountIDs([]AccountID{nodeAccountID}).
			SetAccountID(nodeAccountID).
			Execute(client)
		require.NoError(t, err)
	}
}
