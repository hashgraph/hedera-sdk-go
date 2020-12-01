package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileUpdateTransaction_Execute(t *testing.T) {
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

	var newContents = []byte("Good Night, World")

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	resp, err = NewFileUpdateTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs(nodeIDs).
		SetContents(newContents).
		Execute(client)

	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	contents, err := NewFileContentsQuery().
		SetFileID(fileID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, newContents, contents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
