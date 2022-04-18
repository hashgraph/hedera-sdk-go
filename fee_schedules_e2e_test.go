//go:build all || e2e
// +build all e2e

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

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func DisabledTestIntegrationFeeSchedulesFromBytes(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)

	feeSchedulesBytes, err := NewFileContentsQuery().
		SetFileID(FileID{Shard: 0, Realm: 0, File: 111}).
		Execute(env.Client)
	require.NoError(t, err)
	feeSchedules, err := FeeSchedulesFromBytes(feeSchedulesBytes)
	require.NoError(t, err)
	assert.NotNil(t, feeSchedules)
	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].FeeData.NodeData.Constant, int64(4498129603))
	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].FeeData.ServiceData.Constant, int64(71970073651))
	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].RequestType, RequestTypeCryptoCreate)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func DisabledTestIntegrationNodeAddressBookFromBytes(t *testing.T) { // nolint
	env := NewIntegrationTestEnv(t)

	nodeAddressBookBytes, err := NewFileContentsQuery().
		SetFileID(FileID{Shard: 0, Realm: 0, File: 101}).
		Execute(env.Client)
	require.NoError(t, err)
	nodeAddressbook, err := NodeAddressBookFromBytes(nodeAddressBookBytes)
	require.NoError(t, err)
	assert.NotNil(t, nodeAddressbook)

	for _, ad := range nodeAddressbook.NodeAddresses {
		println(ad.NodeID)
		println(string(ad.CertHash))
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
