//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

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
