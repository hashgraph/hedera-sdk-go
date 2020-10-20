package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReceiptQueryTransaction_Execute(t *testing.T) {
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
		SetNodeAccountID(AccountID{Account: 3}).
		FreezeWith(client)
	assert.NoError(t, err)

	tx.SignWithOperator(client)

	resp, err := tx.Execute(client)

	println("NodeID", resp.NodeID.String())

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	println("AccountID", receipt.AccountID.String())

	record, err := resp.GetRecord(client)
	println("Record", record.TransactionID.String())

	accountID := *record.Receipt.AccountID
	println("AccountID2", accountID)

	assert.NotNil(t, accountID)

	delete, err := NewAccountDeleteTransaction().
		SetNodeAccountID(resp.NodeID).
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetMaxTransactionFee(NewHbar(1)).
		FreezeWith(client)
	assert.NoError(t, err)

	respdelete, err := delete.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = respdelete.GetReceipt(client)
	assert.NoError(t, err)
}
