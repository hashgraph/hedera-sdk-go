package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileCreateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	client.SetMaxTransactionFee(NewHbar(2))
	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}

func Test_FileCreate_NoKey(t *testing.T) {
	client := newTestClient(t)

	client.SetMaxTransactionFee(NewHbar(2))

	resp, err := NewFileCreateTransaction().
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("exceptional precheck status UNAUTHORIZED"), err.Error())
}
