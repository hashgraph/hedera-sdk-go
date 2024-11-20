//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountIDCanPopulateAccountNumber(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.SetIncludeChildren(true).SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	idMirror, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	time.Sleep(5 * time.Second)
	error := idMirror.PopulateAccount(env.Client)
	require.NoError(t, error)
	require.Equal(t, newAccountId.Account, idMirror.Account)
}

func TestIntegrationAccountIDCanPopulateAccountAliasEvmAddress(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.SetIncludeChildren(true).SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	time.Sleep(5 * time.Second)
	error := newAccountId.PopulateEvmAddress(env.Client)
	require.NoError(t, error)
	require.Equal(t, evmAddress, hex.EncodeToString(*newAccountId.AliasEvmAddress))
}

func TestIntegrationAccountIDCanPopulateAccountAliasEvmAddressWithMirror(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.SetIncludeChildren(true).SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	time.Sleep(5 * time.Second)
	error := newAccountId.PopulateEvmAddress(env.Client)
	require.NoError(t, error)
	require.Equal(t, evmAddress, hex.EncodeToString(*newAccountId.AliasEvmAddress))
}

func TestIntegrationAccountIDCanPopulateAccountAliasEvmAddressWithNoMirror(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	env.Client.mirrorNetwork = nil
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.SetIncludeChildren(true).SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	env.Client.mirrorNetwork = nil
	time.Sleep(5 * time.Second)
	error := newAccountId.PopulateEvmAddress(env.Client)
	require.Error(t, error)
}

func TestIntegrationAccountIDCanPopulateAccountAliasEvmAddressWithMirrorAndNoEvmAddress(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	env.Client.mirrorNetwork = nil
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.SetIncludeChildren(true).SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	env.Client.mirrorNetwork = nil
	time.Sleep(5 * time.Second)
	error := newAccountId.PopulateAccount(env.Client)
	require.Error(t, error)
}
