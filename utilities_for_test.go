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
	AccountID:  &accountIDForTransactionID,
	ValidStart: &validStartForTransacionID,
}

type IntegrationTestEnv struct {
	Client              *Client
	OperatorKey         PrivateKey
	OperatorID          AccountID
	OriginalOperatorKey PublicKey
	OriginalOperatorID  AccountID
	NodeAccountIDs      []AccountID
}

func NewIntegrationTestEnv(t *testing.T) IntegrationTestEnv {
	var env IntegrationTestEnv
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		env.Client = ClientForPreviewnet()
	} else if os.Getenv("HEDERA_NETWORK") == "localhost" {
		network := make(map[string]AccountID)
		network["127.0.0.1:50213"] = AccountID{Account: 3}
		network["127.0.0.1:50214"] = AccountID{Account: 4}
		network["127.0.0.1:50215"] = AccountID{Account: 5}

		env.Client = ClientForNetwork(network)
	} else {
		env.Client, err = ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			env.Client = ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		env.OperatorID, err = AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		env.OperatorID.setNetworkWithClient(env.Client)

		env.OperatorKey, err = PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		env.Client.SetOperator(env.OperatorID, env.OperatorKey)
	}

	assert.NotNil(t, env.Client.GetOperatorAccountID())
	assert.NotNil(t, env.Client.GetOperatorPublicKey())

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(NewHbar(30)).
		SetAutoRenewPeriod(time.Hour*24*81 + time.Minute*26 + time.Second*39).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	env.OriginalOperatorID = env.Client.GetOperatorAccountID()
	env.OriginalOperatorKey = env.Client.GetOperatorPublicKey()
	env.OperatorID = *receipt.AccountID
	env.OperatorKey = newKey
	env.NodeAccountIDs = []AccountID{resp.NodeID}
	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	return env
}

func CloseIntegrationTestEnv(env IntegrationTestEnv, token *TokenID) error {
	var resp TransactionResponse
	var err error
	if token != nil {
		dissociateTx, err := NewTokenDeleteTransaction().
			SetNodeAccountIDs(env.NodeAccountIDs).
			SetTokenID(*token).
			FreezeWith(env.Client)
		if err != nil {
			return err
		}

		resp, err = dissociateTx.
			Sign(env.OperatorKey).
			Execute(env.Client)
		if err != nil {
			return err
		}

		_, err = resp.GetReceipt(env.Client)
		if err != nil {
			return err
		}
	}

	resp, err = NewAccountDeleteTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.OperatorID).
		SetTransferAccountID(env.OriginalOperatorID).
		Execute(env.Client)
	if err != nil {
		return err
	}

	_, err = resp.GetReceipt(env.Client)
	if err != nil {
		return err
	}

	return nil
}

func newMockClient() (*Client, error) {
	privateKey, err := PrivateKeyFromString(mockPrivateKey)

	if err != nil {
		return nil, err
	}

	var net = make(map[string]AccountID)
	net["nonexistent-testnet"] = AccountID{Account: 3}

	client := newClient(net, []string{}, "testnet")
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
		SetNodeAccountIDs([]AccountID{AccountID{0, 0, 4, nil}}).
		FreezeWith(client)
	if err != nil {
		return &TransferTransaction{}, err
	}

	tx.Sign(privateKey)

	return tx, nil
}

func newTestClient(t *testing.T, token bool) *Client {
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

	if token {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		resp, err := NewAccountCreateTransaction().
			SetKey(newKey.PublicKey()).
			SetInitialBalance(NewHbar(20)).
			Execute(client)
		assert.NoError(t, err)

		receipt, err := resp.GetReceipt(client)
		assert.NoError(t, err)

		client.SetOperator(*receipt.AccountID, newKey)

		time.Sleep(2000)
	}

	return client
}

func TestTls(t *testing.T) {
	var env IntegrationTestEnv
	var err error

	env.Client = ClientForPreviewnet()

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		env.OperatorID, err = AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		env.OperatorID.setNetworkWithClient(env.Client)

		env.OperatorKey, err = PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		env.Client.SetOperator(env.OperatorID, env.OperatorKey)
	}

	assert.NotNil(t, env.Client.GetOperatorAccountID())
	assert.NotNil(t, env.Client.GetOperatorPublicKey())

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(NewHbar(30)).
		SetAutoRenewPeriod(time.Hour*24*81 + time.Minute*26 + time.Second*39).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)
}
