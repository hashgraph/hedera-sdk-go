//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationLiveHashQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_hash, _ := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	resp2, err := NewLiveHashAddTransaction().
		SetAccountID(accountID).
		SetDuration(24 * 30 * time.Hour).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash).
		SetKeys(newKey.PublicKey()).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED received for transaction %s", resp2.TransactionID), err.Error())
	}

	_, err = NewLiveHashQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		SetHash(_hash).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status NOT_SUPPORTED", err.Error())
	}

	resp2, err = NewLiveHashDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED received for transaction %s", resp2.TransactionID), err.Error())
	}

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationLiveHashQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_hash, _ := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	resp2, err := NewLiveHashAddTransaction().
		SetAccountID(accountID).
		SetDuration(24 * 30 * time.Hour).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash).
		SetKeys(newKey.PublicKey()).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED received for transaction %s", resp2.TransactionID), err.Error())
	}

	liveHashQ := NewLiveHashQuery().
		SetAccountID(accountID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash)

	cost, err := liveHashQ.GetCost(env.Client)
	assert.Error(t, err)

	_, err = liveHashQ.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status NOT_SUPPORTED", err.Error())
	}

	resp2, err = NewLiveHashDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetHash(_hash).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED received for transaction %s", resp2.TransactionID), err.Error())
	}

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}
