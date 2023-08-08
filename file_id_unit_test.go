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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitFileIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := FileIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	id.ToStringWithChecksum(*client)
	sol := id.ToSolidityAddress()
	FileIDFromSolidityAddress(sol)
	id.Validate(client)

	pb := id._ToProtobuf()
	_FileIDFromProtobuf(pb)

	idByte := id.ToBytes()
	FileIDFromBytes(idByte)

	require.Equal(t, FileID{File: 111}.String(), FileIDForFeeSchedule().String())
	require.Equal(t, FileID{File: 102}.String(), FileIDForAddressBook().String())
	require.Equal(t, FileID{File: 112}.String(), FileIDForExchangeRate().String())

	assert.Equal(t, id.File, uint64(123))
}

func TestUnitFileIDChecksumToString(t *testing.T) {
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}
