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
	accountID.checksum = nil;
	return AssessedCustomFee{
		Amount:                100,
		TokenID:               nil,
		FeeCollectorAccountId: &accountID,
		PayerAccountIDs:       []*AccountID{},
	}
}
