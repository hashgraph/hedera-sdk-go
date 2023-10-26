//go:build all || e2e
// +build all e2e

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
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAccountIDCanPopulateAccountNumber(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.GetReceiptQuery().SetIncludeChildren(true).Execute(env.Client)
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

	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.GetReceiptQuery().SetIncludeChildren(true).Execute(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	time.Sleep(5 * time.Second)
	error:= newAccountId.PopulateEvmAddress(env.Client)
	require.NoError(t, error)
	require.Equal(t, evmAddress, hex.EncodeToString(*newAccountId.AliasEvmAddress))
}

func TestIntegrationAccountIDCanPopulateAccountAliasEvmAddressWithMirror(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	tx, err := NewTransferTransaction().AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).Execute(env.Client)
	require.NoError(t, err)
	receipt, err := tx.GetReceiptQuery().SetIncludeChildren(true).Execute(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	time.Sleep(5 * time.Second)
	error:= newAccountId.PopulateEvmAddress(env.Client)
	require.NoError(t, error)
	require.Equal(t, evmAddress, hex.EncodeToString(*newAccountId.AliasEvmAddress))
}

func TestIntegrationAccountIDCanPopulateAccountAliasEvmAddressWithNoMirror(t *testing.T){
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
	receipt, err := tx.GetReceiptQuery().SetIncludeChildren(true).Execute(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	env.Client.mirrorNetwork = nil
	time.Sleep(5 * time.Second)
	error:= newAccountId.PopulateEvmAddress(env.Client)
	require.Error(t, error)
}

func TestIntegrationAccountIDCanPopulateAccountAliasEvmAddressWithMirrorAndNoEvmAddress(t *testing.T){
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
	receipt, err := tx.GetReceiptQuery().SetIncludeChildren(true).Execute(env.Client)
	require.NoError(t, err)
	newAccountId := *receipt.Children[0].AccountID
	env.Client.mirrorNetwork = nil
	time.Sleep(5 * time.Second)
	error:= newAccountId.PopulateAccount(env.Client)
	require.Error(t, error)
}
