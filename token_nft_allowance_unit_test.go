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
)

func TestUnitNewTokenNftAllowance(t *testing.T) {
	tokenID := TokenID{Token: 1}
	owner := AccountID{Account: 2}
	spender := AccountID{Account: 3}
	serialNumbers := []int64{1, 2, 3}
	approvedForAll := true
	delegatingSpender := AccountID{Account: 4}

	allowance := NewTokenNftAllowance(tokenID, owner, spender, serialNumbers, approvedForAll, delegatingSpender)

	assert.Equal(t, &tokenID, allowance.TokenID)
	assert.Equal(t, &owner, allowance.OwnerAccountID)
	assert.Equal(t, &spender, allowance.SpenderAccountID)
	assert.Equal(t, serialNumbers, allowance.SerialNumbers)
	assert.Equal(t, approvedForAll, allowance.AllSerials)
	assert.Equal(t, &delegatingSpender, allowance.DelegatingSpender)
}

func TestUnitTokenNftAllowance_String(t *testing.T) {
	approval := TokenNftAllowance{
		TokenID:           &TokenID{Token: 1},
		SpenderAccountID:  &AccountID{Account: 2},
		OwnerAccountID:    &AccountID{Account: 3},
		SerialNumbers:     []int64{1, 2, 3},
		AllSerials:        true,
		DelegatingSpender: &AccountID{Account: 4},
	}

	assert.Equal(t, "OwnerAccountID: 0.0.3, SpenderAccountID: 0.0.2, TokenID: 0.0.1, Serials: 1, 2, 3, , ApprovedForAll: true", approval.String())
}
