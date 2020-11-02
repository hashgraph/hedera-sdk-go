package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			println("not error", err.Error())
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")
	var operatorKey hedera.PrivateKey

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err = hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	keys := make([]hedera.PrivateKey, 3)
	pubKey := hedera.PublicKey{}

	fmt.Println("threshold key example")
	fmt.Println("Keys: ")

	for i := range keys {
		newKey, err := hedera.GeneratePrivateKey()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Key %v:\n", i)
		fmt.Printf("private = %v\n", newKey)
		fmt.Printf("public = %v\n", newKey.PublicKey())

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	// A threshold key with a threshold of 2 and length of 3 requires
	// at least 2 of the 3 keys to sign anything modifying the account

	//fmt.Printf("threshold key %v\n", thresholdKey)

	transaction, err := hedera.NewAccountCreateTransaction().
		SetKey(pubKeys).
		SetInitialBalance(hedera.NewHbar(6)).
		SetTransactionID(hedera.NewTransactionID(operatorAccountID)).
		SetTransactionMemo("sdk example create_account_with_threshold_keys/main.go").
		Build(client)

	if err != nil {
		panic(err)
	}

	transactionID, err := transaction.Sign(operatorPrivateKey).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	newAccountID := transactionReceipt.GetAccountID()

	fmt.Printf("account = %v\n", newAccountID)

	transferTx, err := hedera.NewCryptoTransferTransaction().
		SetTransactionID(hedera.NewTransactionID(newAccountID)).
		AddSender(newAccountID, hedera.HbarFrom(5, hedera.HbarUnits.Hbar)).
		AddRecipient(operatorAccountID, hedera.HbarFrom(5, hedera.HbarUnits.Hbar)).
		Build(client)

	if err != nil {
		panic(err)
	}

	// Manually sign with 2 of the private keys provided in the threshold
	transferID, err := transferTx.
		Sign(keys[0]).
		Sign(keys[1]).
		Execute(client)

	if err != nil {
		panic(err)
	}

	// Must wait for the transaction to go to consensus
	receipt, err := transferID.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("status of transfer transaction: %v\n", receipt.Status)

	// Operator must be set
	client.SetOperator(operatorAccountID, operatorPrivateKey)

	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(newAccountID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("account balance after transfer: %v\n", balance)
}
