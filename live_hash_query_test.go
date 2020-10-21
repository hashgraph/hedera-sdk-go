package hedera

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestLiveHashQuery_Execute(t *testing.T) {
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
		SetNodeAccountID(AccountID{Account: 3}).
		Execute(client)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewLiveHashAddTransaction().
		SetAccountID(accountID).
		SetDuration(24 * 30 * time.Hour).
		SetNodeAccountID(resp.NodeID).
		SetHash(_hash).
		SetKeys(newKey.PublicKey()).
		Execute(client)

	assert.Error(t, err)

	_, err = NewLiveHashDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountID(resp.NodeID).
		SetHash(_hash).
		Execute(client)
	assert.NoError(t, err)

	_, err = NewLiveHashQuery().
		SetAccountID(accountID).
		SetNodeAccountID(resp.NodeID).
		SetHash(_hash).
		Execute(client)
	assert.Error(t, err)

	_, err = NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetNodeAccountID(resp.NodeID).
		SetTransferAccountID(client.GetOperatorID()).
		Execute(client)
	assert.NoError(t, err)
}
