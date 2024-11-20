//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

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
