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
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" && client.GetOperatorPublicKey().Bytes() == nil {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating an account to delete:")
	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	// first create an account
	transactionResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(hedera.NewHbar(2)).
		SetTransactionMemo("go sdk example delete_account/main.go").
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)
	fmt.Println("deleting created account")

	// To delete an account you must do the following:
	deleteTransaction, err := hedera.NewAccountDeleteTransaction().
		// Set the account to be deleted
		SetAccountID(newAccountID).
		// Set an account to transfer to balance of the deleted account to
		SetTransferAccountID(hedera.AccountID{Account: 3}).
		SetTransactionMemo("go sdk example delete_account/main.go").
		FreezeWith(client)

	if err != nil {
		panic(err)
	}

	// Manually sign the transaction with the private key of the account to be deleted
	deleteTransaction = deleteTransaction.Sign(newKey)

	// Execute the transaction
	deleteTransactionResponse, err := deleteTransaction.Execute(client)

	if err != nil {
		panic(err)
	}

	deleteTransactionReceipt, err := deleteTransactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("account delete transaction status: %v\n", deleteTransactionReceipt.Status)
}
