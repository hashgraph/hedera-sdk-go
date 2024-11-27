//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestSetKeyUsesAnyKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	keys := make([]PrivateKey, 3)
	pubKeys := make([]PublicKey, 3)

	for i := range keys {
		newKey, err := PrivateKeyGenerateEd25519()
		if err != nil {
			panic(err)
		}

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	thresholdKey := KeyListWithThreshold(2).
		AddAllPublicKeys(pubKeys)

	_, err = NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(newKey.PublicKey()).
		SetKey(newKey).
		SetKey(thresholdKey).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	require.NoError(t, err)
}

func DisabledTestECDSAPrivateKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		SetMaxAutomaticTokenAssociations(100).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(accountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		SetTransactionID(TransactionIDGenerate(accountID)).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = tx.
		Sign(newKey).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}
