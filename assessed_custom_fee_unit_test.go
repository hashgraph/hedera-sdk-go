//go:build all || unit
// +build all unit

package hiero

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

// The test checks the conversation methods on the AssessedCustomFee struct. We check wether it is correctly converted to protobuf and back.
func TestUnitassessedCustomFee(t *testing.T) {
	t.Parallel()
	assessedFeeOriginal := _MockAssessedCustomFee()
	assessedFeeBytes := assessedFeeOriginal.ToBytes()
	assessedFeeFromBytes, err := AssessedCustomFeeFromBytes(assessedFeeBytes)
	require.NoError(t, err)
	require.Equal(t, assessedFeeOriginal, assessedFeeFromBytes)
}

func _MockAssessedCustomFee() AssessedCustomFee {
	accountID, _ := AccountIDFromString("0.0.123-esxsf")
	accountID.checksum = nil
	return AssessedCustomFee{
		Amount:                100,
		TokenID:               nil,
		FeeCollectorAccountId: &accountID,
		PayerAccountIDs:       []*AccountID{},
	}
}
