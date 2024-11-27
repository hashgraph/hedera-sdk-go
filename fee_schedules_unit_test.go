//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

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
