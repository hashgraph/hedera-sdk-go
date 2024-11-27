//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationFileInfoQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	info, err := NewFileInfoQuery().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, *fileID, info.FileID)
	assert.Equal(t, info.Size, int64(12))
	assert.False(t, info.IsDeleted)
	assert.NotNil(t, info.Keys)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileInfoQueryGetCost(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	fileInfo := NewFileInfoQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileInfo.GetCost(env.Client)
	require.NoError(t, err)

	info, err := fileInfo.SetQueryPayment(cost).Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, *fileID, info.FileID)
	assert.Equal(t, info.Size, int64(12))
	assert.False(t, info.IsDeleted)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileInfoQuerySetBigMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	fileInfo := NewFileInfoQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(10000)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = fileInfo.GetCost(env.Client)
	require.NoError(t, err)

	info, err := fileInfo.Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, *fileID, info.FileID)
	assert.Equal(t, info.Size, int64(12))
	assert.False(t, info.IsDeleted)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileInfoQuerySetSmallMaxPayment(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	fileInfo := NewFileInfoQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileInfo.GetCost(env.Client)
	require.NoError(t, err)

	_, err = fileInfo.Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "cost of FileInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tℏ", err.Error())
	}

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileInfoQueryInsufficientFee(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	fileInfo := NewFileInfoQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = fileInfo.GetCost(env.Client)
	require.NoError(t, err)

	_, err = fileInfo.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
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

func TestIntegrationFileInfoQueryNoFileID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewFileInfoQuery().
		SetQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional precheck status INVALID_FILE_ID", err.Error())
	}

}
