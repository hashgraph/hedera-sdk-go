package hedera

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestLiveHashQuery_Execute(t *testing.T) {
	var client *Client

	network := os.Getenv("HEDERA_NETWORK")

	if network == "previewnet" {
		client = ClientForPreviewnet()
	}

	client, err := ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	_hash, err := hex.DecodeString("100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {

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

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	resp, err = NewLiveHashAddTransaction().
		SetAccountID(accountID).
		SetDuration(24 * 30 * time.Hour).
		SetNodeAccountIDs(nodeIDs).
		SetHash(_hash).
		SetKeys(newKey.PublicKey()).
		Execute(client)

	assert.Error(t, err)

	_, err = NewLiveHashDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetHash(_hash).
		Execute(client)
	assert.NoError(t, err)

	_, err = NewLiveHashQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetHash(_hash).
		Execute(client)
	assert.Error(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountIDs(nodeIDs).
		SetTransferAccountID(client.GetOperatorID()).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
