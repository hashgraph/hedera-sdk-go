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

func TestUnitNewTokenAllowance(t *testing.T) {
	t.Parallel()

	tokID := TokenID{Token: 3}
	owner := AccountID{Account: 5}
	spender := AccountID{Account: 6}
	amount := int64(100)

	allowance := NewTokenAllowance(tokID, owner, spender, amount)

	newAllowance := TokenAllowance{
		TokenID:          &tokID,
		SpenderAccountID: &spender,
		OwnerAccountID:   &owner,
		Amount:           amount,
	}

	assert.Equal(t, newAllowance, allowance)
}

func TestUnitTokenAllowanceFromProtobuf(t *testing.T) {
	t.Parallel()

	tokID := TokenID{Token: 3}
	owner := AccountID{Account: 5}
	spender := AccountID{Account: 6}
	amount := int64(100)

	allowance := NewTokenAllowance(tokID, owner, spender, amount)

	pb := allowance._ToProtobuf()
	assert.NotNil(t, pb)

	allowance2 := _TokenAllowanceFromProtobuf(pb)
	assert.Equal(t, allowance, allowance2)
}

func TestUnitTokenAllowance_String(t *testing.T) {
	t.Parallel()

	tokID := TokenID{Token: 3}
	owner := AccountID{Account: 5}
	spender := AccountID{Account: 6}
	amount := int64(100)

	allowance := NewTokenAllowance(tokID, owner, spender, amount)

	assert.Equal(t, "OwnerAccountID: 0.0.5, SpenderAccountID: 0.0.6, TokenID: 0.0.3, Amount: 100 tℏ", allowance.String())
}
