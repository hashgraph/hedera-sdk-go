package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	client.SetOperator(operatorAccountID, operatorKey)

	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey")
		return
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
		println(err.Error(), ": error creating account")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account creation receipt")
		return
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
		println(err.Error(), ": error freezing account delete transaction")
		return
	}

	// Manually sign the transaction with the private key of the account to be deleted
	deleteTransaction = deleteTransaction.Sign(newKey)

	// Execute the transaction
	deleteTransactionResponse, err := deleteTransaction.Execute(client)

	if err != nil {
		println(err.Error(), ": error deleting account")
		return
	}

	deleteTransactionReceipt, err := deleteTransactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving account deletion receipt")
		return
	}

	fmt.Printf("account delete transaction status: %v\n", deleteTransactionReceipt.Status)
}
