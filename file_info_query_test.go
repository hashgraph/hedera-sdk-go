package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestNewFileInfoQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewFileInfoQuery().
		SetFileID(FileID{File: 3}).
		SetQueryPaymentTransaction(mockTransaction)

	assert.Equal(t, `fileGetInfo:{header:{payment:{bodyBytes:"\n\x0e\n\x08\x08\xdc\xc9\x07\x10۟\t\x12\x02\x18\x03\x12\x02\x18\x03\x18\x80\xc2\xd7/\"\x02\x08xr\x14\n\x12\n\x07\n\x02\x18\x02\x10\xc7\x01\n\x07\n\x02\x18\x03\x10\xc8\x01"sigMap:{sigPair:{pubKeyPrefix:"\xe4\xf1\xc0\xebL}\xcd\xc3\xe7\xeb\x11p\xb3\x08\x8a=\x12\xa2\x97\xf4\xa3\xeb\xe2\xf2\x85\x03\xfdg5F\xed\x8e"ed25519:"\x12&5\x96\xfb\xb4\x1c]P\xbb%\xecP\x9bk͙\x0b߼\xac)\xa6+\xd2<\x97+\xbb\x8c\x8af\xcb\xdai\x17T4{\xf7\xf3UYn\n\x8f\xabep\x04\xf6\x83\x0f\xbaFUP\xa3\xd1/\x1d\x9d\x1a\x0b"}}}}fileID:{fileNum:3}}`, strings.ReplaceAll(query.QueryBuilder.pb.String(), " ", ""))
}

func TestFileInfoQueryTransaction_Execute(t *testing.T) {
	client, err := ClientFromFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := Ed25519PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	client.SetMaxTransactionFee(NewHbar(2))

	txID, err := NewFileCreateTransaction().
		AddKey(client.GetOperatorKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk e2e tests").
		Execute(client)

	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.fileID
	assert.NotNil(t, fileID)

	info, err := NewFileInfoQuery().
		SetFileID(*fileID).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, *fileID, info.FileID)

	txID, err = NewFileDeleteTransaction().
		SetFileID(*fileID).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
