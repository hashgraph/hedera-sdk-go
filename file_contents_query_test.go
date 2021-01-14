package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeFileContentsQuery(t *testing.T) {
	query := NewFileContentsQuery().
		SetFileID(FileID{File: 3}).
		Query

	assert.Equal(t, `fileGetContents:{header:{}fileID:{fileNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestFileContentsQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	remoteContents, err := NewFileContentsQuery().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func TestFileContentsQueryCost_Execute(t *testing.T) {
	client := newTestClient(t)

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(contents).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := fileContents.GetCost(client)
	assert.NoError(t, err)

	remoteContents, err := fileContents.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_FileContents_NoFileID(t *testing.T) {
	client := newTestClient(t)

	_, err := NewFileContentsQuery().
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_FILE_ID"), err.Error())
	}
}
