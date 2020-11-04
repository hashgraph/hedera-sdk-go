package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var mockPrivateKey string = "302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962"

var testTransactionID TransactionID = TransactionID{
	AccountID{Account: 3},
	time.Unix(124124, 151515),
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

func newMockTransaction() (Transaction, error) {
	privateKey, err := PrivateKeyFromString(mockPrivateKey)

	if err != nil {
		return Transaction{}, err
	}

	client, err := newMockClient()

	if err != nil {
		return Transaction{}, err
	}

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0], err = AccountIDFromString("0.0.4")
	if err != nil {
		return Transaction{}, err
	}

	tx, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		SetMaxTransactionFee(HbarFrom(1, HbarUnits.Hbar)).
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs(nodeIDs).
		FreezeWith(client)

	if err != nil {
		return Transaction{}, err
	}

	tx.Sign(privateKey)

	return tx.Transaction, nil
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
