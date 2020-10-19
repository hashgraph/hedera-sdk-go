package hedera

import (
	"os"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountBalanceQuery(t *testing.T) {
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


	response, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := response.TransactionID.GetReceipt(client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	query, err := NewAccountBalanceQuery().
		SetAccountID(account).
		SetNodeAccountID(response.NodeID).
		Execute(client)
	assert.NoError(t, err)

	println("Success", query.String())
}

func TestNewAccountBalanceQuery_ForContract(t *testing.T) {
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


	response, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(NewHbar(1)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := response.TransactionID.GetReceipt(client)
	assert.NoError(t, err)

	contract := *receipt.ContractID

	query, err := NewAccountBalanceQuery().
		SetContractID(contract).
		SetNodeAccountID(response.NodeID).
		Execute(client)
	assert.NoError(t, err)

	println("Success", query.String())
}

func TestAccountBalanceQuery_Execute(t *testing.T) {
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

	response, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := response.GetReceipt(client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID
	assert.NoError(t, err)

	balance, err := NewAccountBalanceQuery().
		SetAccountID(accountID).
		Execute(client)
	assert.NoError(t, err)
	assert.Equal(t, newBalance, balance)

	tx, err := NewAccountDeleteTransaction().
		SetAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(NewTransactionIDWithValidStart(accountID, time.Now())).
		FreezeWith(client)
	assert.NoError(t, err)

	response, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = response.GetReceipt(client)
	assert.NoError(t, err)
}
