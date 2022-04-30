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

func TestUnitTransactionID(t *testing.T) {
	txID := TransactionIDGenerate(AccountID{0, 0, 3, nil, nil, nil})
	txID = txID.SetScheduled(true)
}

func TestUnitTransactionIDFromString(t *testing.T) {
	txID, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled")
	require.NoError(t, err)
	require.Equal(t, txID.AccountID.String(), "0.0.3")
	require.True(t, txID.scheduled)
}

func TestUnitTransactionIDFromStringNonce(t *testing.T) {
	txID, err := TransactionIdFromString("0.0.3@1614997926.774912965?scheduled/4")
	require.NoError(t, err)
	require.Equal(t, *txID.Nonce, int32(4))
	require.Equal(t, txID.AccountID.String(), "0.0.3")
}
