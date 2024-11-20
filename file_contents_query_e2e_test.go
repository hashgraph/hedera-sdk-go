//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationFileContentsQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	remoteContents, err := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileContentsQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileContents.GetCost(env.Client)
	require.NoError(t, err)

	remoteContents, err := fileContents.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileContentsQuerySetBigMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(100000)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileContents.GetCost(env.Client)
	require.NoError(t, err)

	remoteContents, err := fileContents.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileContentsQuerySetSmallMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileContents.GetCost(env.Client)
	require.NoError(t, err)

	_, err = fileContents.Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "cost of FileContentsQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tℏ", err.Error())
	}

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileContentsQueryInsufficientFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = fileContents.GetCost(env.Client)
	require.NoError(t, err)

	_, err = fileContents.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INSUFFICIENT_TX_FEE", err.Error())
	}

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileContentsQueryNoFileID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewFileContentsQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INVALID_FILE_ID", err.Error())
	}

}
