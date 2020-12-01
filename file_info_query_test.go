package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeFileInfoQuery(t *testing.T) {
	query := NewFileInfoQuery().
		SetFileID(FileID{File: 3}).
		Query

	assert.Equal(t, `fileGetInfo:{header:{}fileID:{fileNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestFileInfoQueryTransaction_Execute(t *testing.T) {
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

	fileID := receipt.FileID
	assert.NotNil(t, fileID)

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	info, err := NewFileInfoQuery().
		SetFileID(*fileID).
		SetNodeAccountIDs(nodeIDs).
		SetQueryPayment(NewHbar(22)).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, *fileID, info.FileID)
	assert.Equal(t, info.Size, int64(12))
	assert.False(t, info.IsDeleted)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
