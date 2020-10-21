package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileContentsQuery_Execute(t *testing.T) {
	client, err := ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	client.SetMaxTransactionFee(NewHbar(2))

	var contents = []byte("Hellow world!")

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorKey()).
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

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	remoteContents, err := NewFileContentsQuery().
		SetFileID(*fileID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, contents, remoteContents)

	resp, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
