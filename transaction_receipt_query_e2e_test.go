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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationTransactionReceiptQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar._NumberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(newBalance).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err = tx.SignWithOperator(env.Client)
	require.NoError(t, err)

	resp, err := tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = resp.GetRecord(env.Client)
	require.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

// what is this
func DisabledTestIntegrationTransactionReceiptQueryInvalidTransactionID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	key, err := GeneratePrivateKey()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(key).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountID := *receipt.AccountID

	tx := NewTransferTransaction().
		AddHbarTransfer(accountID, NewHbar(1).Negated()).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(1))

	scheduleTx, err := tx.Schedule()
	require.NoError(t, err)

	scheduleTx = scheduleTx.
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetPayerAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID()))

	resp, err = scheduleTx.Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	scheduleID := *receipt.ScheduleID

	scheduleSignTx, err := NewScheduleSignTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err = scheduleSignTx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	assert.Error(t, err)

	switch receiptErr := err.(type) {
	case ErrHederaReceiptStatus:
		assert.NotNil(t, receiptErr.Receipt.ExchangeRate)
	default:
		panic("err was not a `ErrHederaReceiptStatus")
	}

	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}
