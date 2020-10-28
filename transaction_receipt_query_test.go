package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReceiptQueryTransaction_Execute(t *testing.T) {
	var client *Client

	network := os.Getenv("HEDERA_NETWORK")

	if network == "previewnet" {
		client = ClientForPreviewnet()
	}

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

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		FreezeWith(client)
	assert.NoError(t, err)

	tx.SignWithOperator(client)

	resp, err := tx.Execute(client)

	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	record, err := resp.GetRecord(client)
	assert.NoError(t, err)

	accountID := *record.Receipt.AccountID
	assert.NotNil(t, accountID)

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	transcation, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs(nodeIDs).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)
	assert.NoError(t, err)

	resp, err = transcation.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
