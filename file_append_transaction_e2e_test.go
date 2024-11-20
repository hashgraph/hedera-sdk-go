//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestIntegrationFileAppendTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte("Hello")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	resp, err = NewFileAppendTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetContents([]byte(" world!")).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	contents, err := NewFileContentsQuery().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, []byte("Hello world!"), contents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileAppendTransactionNoFileID(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte("Hello")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	_, err = NewFileAppendTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetContents([]byte(" world!")).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "exceptional precheck status INVALID_FILE_ID", err.Error())
	}

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationFileAppendTransactionNothingSet(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewFileAppendTransaction().
		SetContents([]byte(" world!")).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "exceptional precheck status INVALID_FILE_ID", err.Error())
	}

}
func TestIntegrationFileAppendTransactionCanExecuteAfterSerializationDeserialization(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte("Hello")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	before := NewFileAppendTransaction().
		SetFileID(fileID)

	bytes, err := before.ToBytes()
	require.NoError(t, err)

	afterI, err := TransactionFromBytes(bytes)
	require.NoError(t, err)

	tx := afterI.(FileAppendTransaction)

	resp, err = tx.SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetContents([]byte(" world!")).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	contents, err := NewFileContentsQuery().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	assert.Equal(t, []byte("Hello world!"), contents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}
