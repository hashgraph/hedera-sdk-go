//go:build all || unit
// +build all unit

package hedera

import (
	"fmt"
	"testing"
	"time"

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

func TestUnitAccountAllowanceAdjustTransactionGet(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyGenerate()
	require.NoError(t, err)
	key2, err := PrivateKeyGenerate()
	require.NoError(t, err)

	nodeAccountIDs := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 123})
	tokenID := TokenID{Token: 3}
	NftID := NftID{tokenID, 1}
	tx, err := NewAccountAllowanceAdjustTransaction().
		AddHbarAllowance(nodeAccountIDs[0], HbarFromTinybar(1)).
		AddTokenAllowance(tokenID, nodeAccountIDs[0], 1).
		AddTokenNftAllowance(NftID, nodeAccountIDs[0]).
		SetTransactionID(transactionID).SetNodeAccountIDs(nodeAccountIDs).
		SetMaxTransactionFee(HbarFromTinybar(100)).SetRegenerateTransactionID(true).
		SetTransactionMemo("go sdk unit test").SetTransactionValidDuration(time.Second * 120).
		SetMaxRetry(1).SetMaxBackoff(time.Second * 120).SetMinBackoff(time.Second * 1).
		Freeze()
	sign, err := key2.SignTransaction(&tx.Transaction)
	require.NoError(t, err)
	tx.AddSignature(key.PublicKey(), sign)
	tx.AddSignature(key2.PublicKey(), sign)

	expectedHbarAllowances := []*HbarAllowance{
		{
			SpenderAccountID: &nodeAccountIDs[0],
			OwnerAccountID:   nil,
			Amount:           1,
		},
	}

	expectedTokenAllowances := []*TokenAllowance{
		{
			TokenID:          &tokenID,
			SpenderAccountID: &nodeAccountIDs[0],
			OwnerAccountID:   nil,
			Amount:           1,
		},
	}

	expectedTokenNftAllowances := []*TokenNftAllowance{
		{
			TokenID:          &tokenID,
			SpenderAccountID: &nodeAccountIDs[0],
			OwnerAccountID:   nil,
			SerialNumbers:    []int64{1},
			AllSerials:       false,
		},
	}

	require.NoError(t, err)
	require.Equal(t, expectedHbarAllowances, tx.GetHbarAllowances())
	require.Equal(t, expectedTokenAllowances, tx.GetTokenAllowances())
	require.Equal(t, expectedTokenNftAllowances, tx.GetTokenNftAllowances())
	require.Equal(t, transactionID, tx.GetTransactionID())
	require.Equal(t, nodeAccountIDs, tx.GetNodeAccountIDs())
	require.Equal(t, HbarFromTinybar(100), tx.GetMaxTransactionFee())
	require.Equal(t, true, tx.GetRegenerateTransactionID())
	require.Equal(t, "go sdk unit test", tx.GetTransactionMemo())
	require.Equal(t, time.Second*120, tx.GetTransactionValidDuration())
	require.Equal(t, 1, tx.GetMaxRetry())
	require.Equal(t, time.Second*120, tx.GetMaxBackoff())
	require.Equal(t, time.Second*1, tx.GetMinBackoff())
	require.Equal(t, fmt.Sprintf("AccountAllowanceAdjustTransaction:%v", transactionID.ValidStart.UnixNano()), tx._GetLogID())
}

func TestUnitAccountAllowanceAdjustTransactionGrantHbarAllowance(t *testing.T) {
	t.Parallel()

	tx := NewAccountAllowanceAdjustTransaction().
		GrantHbarAllowance(AccountID{Account: 3}, AccountID{Account: 4}, HbarFromTinybar(1))
	expectedHbarAllowances := []*HbarAllowance{
		{
			SpenderAccountID: &AccountID{Account: 4},
			OwnerAccountID:   &AccountID{Account: 3},
			Amount:           1,
		},
	}
	require.Equal(t, expectedHbarAllowances, tx.GetHbarAllowances())
}
func TestUnitAccountAllowanceAdjustTransactionRevokeHbarAllowance(t *testing.T) {
	t.Parallel()

	tx := NewAccountAllowanceAdjustTransaction().
		RevokeHbarAllowance(AccountID{Account: 3}, AccountID{Account: 4}, HbarFromTinybar(1))
	expectedHbarAllowances := []*HbarAllowance{
		{
			SpenderAccountID: &AccountID{Account: 4},
			OwnerAccountID:   &AccountID{Account: 3},
			Amount:           -1,
		},
	}
	require.Equal(t, expectedHbarAllowances, tx.GetHbarAllowances())
}

func TestUnitAccountAllowanceAdjustTransactionGrantTokenAllowance(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 3}
	tx := NewAccountAllowanceAdjustTransaction().
		GrantTokenAllowance(tokenID, AccountID{Account: 3}, AccountID{Account: 4}, 1)
	expectedTokenAllowances := []*TokenAllowance{
		{
			TokenID:          &tokenID,
			SpenderAccountID: &AccountID{Account: 4},
			OwnerAccountID:   &AccountID{Account: 3},
			Amount:           1,
		},
	}
	require.Equal(t, expectedTokenAllowances, tx.GetTokenAllowances())
}

func TestUnitAccountAllowanceAdjustTransactionRevokeTokenAllowance(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 3}
	tx := NewAccountAllowanceAdjustTransaction().
		RevokeTokenAllowance(tokenID, AccountID{Account: 3}, AccountID{Account: 4}, 1)
	expectedTokenAllowances := []*TokenAllowance{
		{
			TokenID:          &tokenID,
			SpenderAccountID: &AccountID{Account: 4},
			OwnerAccountID:   &AccountID{Account: 3},
			Amount:           -1,
		},
	}
	require.Equal(t, expectedTokenAllowances, tx.GetTokenAllowances())
}

func TestUnitAccountAllowanceAdjustTransactionGrantTokenNftAllowance(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 3}
	NftID := NftID{tokenID, 1}
	tx := NewAccountAllowanceAdjustTransaction().
		GrantTokenNftAllowance(NftID, AccountID{Account: 3}, AccountID{Account: 4})
	expectedTokenNftAllowances := []*TokenNftAllowance{
		{
			SpenderAccountID:  &AccountID{Account: 4},
			OwnerAccountID:    &AccountID{Account: 3},
			TokenID:           &tokenID,
			SerialNumbers:     []int64{1},
			AllSerials:        false,
			DelegatingSpender: nil,
		},
	}
	require.Equal(t, expectedTokenNftAllowances, tx.GetTokenNftAllowances())
}

func TestUnitAccountAllowanceAdjustTransactionRevokeTokenNftAllowance(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 3}
	NftID := NftID{tokenID, 1}
	tx := NewAccountAllowanceAdjustTransaction().
		RevokeTokenNftAllowance(NftID, AccountID{Account: 3}, AccountID{Account: 4})
	expectedTokenNftAllowances := []*TokenNftAllowance{
		{
			SpenderAccountID:  &AccountID{Account: 4},
			OwnerAccountID:    &AccountID{Account: 3},
			TokenID:           &tokenID,
			SerialNumbers:     []int64{1},
			AllSerials:        false,
			DelegatingSpender: nil,
		},
	}
	require.Equal(t, expectedTokenNftAllowances, tx.GetTokenNftAllowances())
}

func TestUnitAccountAllowanceAdjustTransactionAddAllTokenNftAllowance(t *testing.T) {
	t.Parallel()

	tokenID := TokenID{Token: 3}
	tx := NewAccountAllowanceAdjustTransaction().
		AddAllTokenNftAllowance(tokenID, AccountID{Account: 3})
	expectedTokenNftAllowances := []*TokenNftAllowance{
		{
			SpenderAccountID:  &AccountID{Account: 3},
			OwnerAccountID:    nil,
			TokenID:           &tokenID,
			SerialNumbers:     []int64{},
			AllSerials:        true,
			DelegatingSpender: nil,
		},
	}
	require.Equal(t, expectedTokenNftAllowances, tx.GetTokenNftAllowances())
}
