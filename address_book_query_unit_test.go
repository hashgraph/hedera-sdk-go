//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/require"
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

func TestUnitAddressBookQueryCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	file := FileID{File: 3, checksum: &checksum}

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewAddressBookQuery().
		SetFileID(file).
		SetLimit(3).
		SetMaxAttempts(4)

	err = query._ValidateNetworkOnIDs(client)

	require.NoError(t, err)
	query.GetFileID()
	query.GetLimit()
	query.GetFileID()
}
