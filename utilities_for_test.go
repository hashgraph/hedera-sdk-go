package hedera

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	if os.Getenv("HEDERA_NETWORK") == "previewnet" { // nolint
		env.Client = ClientForPreviewnet()
	} else if os.Getenv("HEDERA_NETWORK") == "localhost" {
		network := make(map[string]AccountID)
		network["127.0.0.1:50213"] = AccountID{Account: 3}
		network["127.0.0.1:50214"] = AccountID{Account: 4}
		network["127.0.0.1:50215"] = AccountID{Account: 5}

		env.Client = ClientForNetwork(network)
	} else if os.Getenv("HEDERA_NETWORK") == "testnet" {
		env.Client = ClientForTestnet()
	} else if os.Getenv("CONFIG_FILE") != "" {
		env.Client, err = ClientFromConfigFile(os.Getenv("CONFIG_FILE"))
		if err != nil {
			panic(err)
		}
	} else {
		panic("Failed to construct client from environment variables")
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
		SetNodeAccountIDs([]AccountID{{0, 0, 4, nil}}).
		FreezeWith(client)
	if err != nil {
		return &TransferTransaction{}, err
	}

	tx.Sign(privateKey)

	return tx, nil
}

func TestIntegrationPreviewnetTls(t *testing.T) {
	var network = map[string]AccountID{
		"0.previewnet.hedera.com:50212": {Account: 3},
		"1.previewnet.hedera.com:50212": {Account: 4},
		"2.previewnet.hedera.com:50212": {Account: 5},
		"3.previewnet.hedera.com:50212": {Account: 6},
		"4.previewnet.hedera.com:50212": {Account: 7},
	}

	client := ClientForNetwork(network)

	for _, nodeAccountID := range network {
		_, err := NewAccountBalanceQuery().
			SetNodeAccountIDs([]AccountID{nodeAccountID}).
			SetAccountID(nodeAccountID).
			Execute(client)
		assert.NoError(t, err)
	}
}

func TestIntegrationTestnetTls(t *testing.T) {
	var network = map[string]AccountID{
		"0.testnet.hedera.com:50212": {Account: 3},
		"1.testnet.hedera.com:50212": {Account: 4},
		"2.testnet.hedera.com:50212": {Account: 5},
		"3.testnet.hedera.com:50212": {Account: 6},
		"4.testnet.hedera.com:50212": {Account: 7},
	}

	client := ClientForNetwork(network)

	for _, nodeAccountID := range network {
		_, err := NewAccountBalanceQuery().
			SetNodeAccountIDs([]AccountID{nodeAccountID}).
			SetAccountID(nodeAccountID).
			Execute(client)
		assert.NoError(t, err)
	}
}
