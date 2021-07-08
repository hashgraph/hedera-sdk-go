package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileContentsQuery_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	remoteContents, err := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestFileContentsQueryCost_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileContents.GetCost(env.Client)
	assert.NoError(t, err)

	remoteContents, err := fileContents.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestFileContentsQueryCost_BigMax_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(100000)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileContents.GetCost(env.Client)
	assert.NoError(t, err)

	remoteContents, err := fileContents.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestFileContentsQueryCost_SmallMax_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileContents.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = fileContents.Execute(env.Client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("cost of FileContentsQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tħ"), err.Error())
	}

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestFileContentsQueryCost_InsufficientFee_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = fileContents.GetCost(env.Client)
	assert.NoError(t, err)

	_, err = fileContents.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INSUFFICIENT_TX_FEE"), err.Error())
	}

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func Test_FileContents_NoFileID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewFileContentsQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_FILE_ID"), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
