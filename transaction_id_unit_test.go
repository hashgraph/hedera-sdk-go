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
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitTransactionID(t *testing.T) {
	t.Parallel()

	txID := TransactionIDGenerate(AccountID{0, 0, 3, nil, nil, nil})
	txID = txID.SetScheduled(true)
}

func TestUnitTransactionIDFromString(t *testing.T) {
	t.Parallel()

	txID, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled")
	require.NoError(t, err)
	require.Equal(t, txID.AccountID.String(), "0.0.3")
	require.True(t, txID.scheduled)
}

func TestUnitTransactionIDFromStringNonce(t *testing.T) {
	t.Parallel()

	txID, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled/4")
	require.NoError(t, err)
	require.Equal(t, *txID.Nonce, int32(4))
	require.Equal(t, txID.AccountID.String(), "0.0.3")
}

func TestUnitTransactionIDFromStringLeadingZero(t *testing.T) {
	t.Parallel()

	txID, err := TransactionIdFromString("0.0.3@1614997926.074912965")
	require.NoError(t, err)
	require.Equal(t, txID.String(), "0.0.3@1614997926.074912965")
}

func TestUnitTransactionIDFromStringTrimmedZeroes(t *testing.T) {
	t.Parallel()

	txID, err := TransactionIdFromString("0.0.3@1614997926.5")
	require.NoError(t, err)
	require.Equal(t, txID.String(), "0.0.3@1614997926.000000005")
}

func TestUnitConcurrentTransactionIDsAreUnique(t *testing.T) {
	const numOfTxns = 100000

	account := AccountID{Account: 1}

	// Channel to collect generated transaction IDs
	idsCh := make(chan TransactionID, numOfTxns)

	var wg sync.WaitGroup
	for i := 0; i < numOfTxns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			idsCh <- TransactionIDGenerate(account)
		}()
	}

	// Close idsCh after all goroutines complete
	go func() {
		wg.Wait()
		close(idsCh)
	}()

	seen := make(map[TransactionID]bool)
	for id := range idsCh {
		require.False(t, seen[id], "Transaction ID %v is not unique", id)
		seen[id] = true
	}

	require.Equal(t, len(seen), numOfTxns)
}
