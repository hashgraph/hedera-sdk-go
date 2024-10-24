//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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

func TestUnitTransferTransactionSetTokenTransferWithDecimals(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 1}
	senderAccountID := AccountID{Account: 2}
	amount := int64(10)
	decimals := uint32(5)

	transaction := NewTransferTransaction().
		AddTokenTransferWithDecimals(tokenID, senderAccountID, amount, decimals)

	require.Equal(t, transaction.GetTokenIDDecimals()[tokenID], decimals)
}

func TestUnitTransferTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	transfer := NewTransferTransaction().
		AddHbarTransfer(accountID, HbarFromTinybar(1))

	err = transfer.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransferTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	transfer := NewTransferTransaction().
		AddHbarTransfer(accountID, HbarFromTinybar(1))

	err = transfer.validateNetworkOnIDs(client)
	require.Error(t, err)
}
