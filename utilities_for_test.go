package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var mockPrivateKey string = "302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962"

var accountIDForTransactionID = AccountID{Account: 3}
var validStartForTransacionID = time.Unix(124124, 151515)

var testTransactionID TransactionID = TransactionID{
    AccountID: &accountIDForTransactionID,
    ValidStart: &validStartForTransacionID, 
}

func newMockClient() (*Client, error) {
	privateKey, err := PrivateKeyFromString(mockPrivateKey)

	if err != nil {
		return nil, err
	}

	var net = make(map[string]AccountID)
	net["nonexistent-testnet"] = AccountID{Account: 3}

	client := newClient(net, []string{})
	client.SetOperator(AccountID{Account: 2}, privateKey)

	return client, nil
}

func newMockTransaction() (*TransferTransaction, error) {
	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	if err != nil {
		return &TransferTransaction{}, err
	}

	client, err := newMockClient()
	if err != nil {
		return &TransferTransaction{}, err
	}

	tx, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs([]AccountID{AccountID{0, 0, 4}}).
		FreezeWith(client)
	if err != nil {
		return &TransferTransaction{}, err
	}

	tx.Sign(privateKey)

	return tx, nil
}

func newTestClient(t *testing.T) *Client {
	var client *Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = ClientForPreviewnet()
	} else {
		client, err = ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			client = ClientForTestnet()
		}

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

	return client
}
