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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitFeeSchedulesFromBytes(t *testing.T) {
	t.Parallel()
	// nolint
	dat, err := os.ReadFile("./fee_schedule/fee_schedule.pb")
	require.NoError(t, err)
	feeSchedules, err := FeeSchedulesFromBytes(dat)
	require.NoError(t, err)
	assert.NotNil(t, feeSchedules)
	assert.Equal(t, int64(11461413665), feeSchedules.current.TransactionFeeSchedules[0].Fees[0].NodeData.Constant)
	assert.Equal(t, int64(229228273302), feeSchedules.current.TransactionFeeSchedules[0].Fees[0].ServiceData.Constant)
	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].RequestType, RequestTypeCryptoCreate)
}
