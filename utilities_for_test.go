//go:build all || unit || e2e || testnets
// +build all unit e2e testnets

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var mockPrivateKey string = "302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962"

var accountIDForTransactionID = AccountID{Account: 3}
var validStartForTransacionID = time.Unix(124124, 151515)

var testTransactionID TransactionID = TransactionID{
	AccountID:  &accountIDForTransactionID,
	ValidStart: &validStartForTransacionID,
}

const testClientJSON string = `{
    "network": {
		"35.237.200.180:50211": "0.0.3",
		"35.186.191.247:50211": "0.0.4",
		"35.192.2.25:50211": "0.0.5",
		"35.199.161.108:50211": "0.0.6",
		"35.203.82.240:50211": "0.0.7",
		"35.236.5.219:50211": "0.0.8",
		"35.197.192.225:50211": "0.0.9",
		"35.242.233.154:50211": "0.0.10",
		"35.240.118.96:50211": "0.0.11",
		"35.204.86.32:50211": "0.0.12"
    },
    "mirrorNetwork": "testnet"
}`

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
		mirror := []string{"127.0.0.1:5600"}
		env.Client = ClientForNetwork(network)
		env.Client.SetMirrorNetwork(mirror)
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
		require.NoError(t, err)

		env.OperatorKey, err = PrivateKeyFromString(configOperatorKey)
		require.NoError(t, err)

		env.Client.SetOperator(env.OperatorID, env.OperatorKey)
	}

	assert.NotNil(t, env.Client.GetOperatorAccountID())
	assert.NotNil(t, env.Client.GetOperatorPublicKey())

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	env.Client.SetMaxNodeAttempts(1)
	env.Client.SetMinBackoff(250 * time.Millisecond)
	env.Client.SetMaxBackoff(8 * time.Second)
	env.Client.SetNodeMinReadmitPeriod(5 * time.Second)
	env.Client.SetNodeMaxReadmitPeriod(1 * time.Hour)
	env.Client.SetMaxAttempts(100)
	env.Client.SetDefaultMaxTransactionFee(NewHbar(50))
	env.Client.SetDefaultMaxQueryPayment(NewHbar(50))

	network := make(map[string]AccountID)

	for key, value := range env.Client.GetNetwork() {
		_, err = NewAccountBalanceQuery().
			SetNodeAccountIDs([]AccountID{value}).
			SetAccountID(value).
			Execute(env.Client)

		if err != nil {
			println(err.Error())
			continue
		}

		network[key] = value
		break
	}

	_ = env.Client.SetNetwork(network)

	if len(network) == 0 {
		panic("failed to construct network; each node returned an error")
	}

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(NewHbar(150)).
		SetAutoRenewPeriod(time.Hour*24*81 + time.Minute*26 + time.Second*39).
		Execute(env.Client)
	if err != nil {
		panic(err)
	}

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	if err != nil {
		panic(err)
	}

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
		deleteTokenTx, err := NewTokenDeleteTransaction().
			SetNodeAccountIDs(env.NodeAccountIDs).
			SetTokenID(*token).
			FreezeWith(env.Client)
		if err != nil {
			return err
		}

		resp, err = deleteTokenTx.
			Sign(env.OperatorKey).
			Execute(env.Client)
		if err != nil {
			return err
		}

		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		if err != nil {
			return err
		}

		// This is needed, because we can't delete the account while still having tokens.
		// This works only, because the token is deleted, otherwise the acount would need to have 0 balance of it before dissociating.
		dissociateTx, err := NewTokenDissociateTransaction().
			SetAccountID(env.Client.operator.accountID).
			SetNodeAccountIDs(env.NodeAccountIDs).
			AddTokenID(*token).
			Execute(env.Client)
		if err != nil {
			return err
		}

		_, err = dissociateTx.SetValidateStatus(true).GetReceipt(env.Client)
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

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	if err != nil {
		return err
	}

	return env.Client.Close()
}

func _NewMockClient() (*Client, error) {
	privateKey, err := PrivateKeyFromString(mockPrivateKey)

	if err != nil {
		return nil, err
	}

	var net = make(map[string]AccountID)
	net["nonexistent-testnet:56747"] = AccountID{Account: 3}

	client := ClientForNetwork(net)
	defaultNetwork := []string{"nonexistent-mirror-testnet:443"}
	client.SetMirrorNetwork(defaultNetwork)
	client.SetOperator(AccountID{Account: 2}, privateKey)

	return client, nil
}

func _NewMockTransaction() (*TransferTransaction, error) {
	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	if err != nil {
		return &TransferTransaction{}, err
	}

	client, err := _NewMockClient()
	if err != nil {
		return &TransferTransaction{}, err
	}

	tx, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs([]AccountID{{0, 0, 4, nil, nil, nil}}).
		FreezeWith(client)
	if err != nil {
		return &TransferTransaction{}, err
	}

	tx.Sign(privateKey)

	return tx, nil
}
