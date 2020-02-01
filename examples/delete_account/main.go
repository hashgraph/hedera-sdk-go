package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	// first create an account
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	newKey, err := hedera.GenerateEd25519PrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating an account to delete:")
	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	client := hedera.ClientForTestnet().
		SetOperator(operatorAccountID, operatorPrivateKey)

	transactionID, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(hedera.NewHbar(2)).
		SetTransactionMemo("go sdk example delete_account/main.go").
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
	fmt.Println("deleting created account")

	// To delete an account you must do the following:
	deleteTransaction, err := hedera.NewAccountDeleteTransaction().
		// Set the account to be deleted
		SetDeleteAccountID(newAccountID).
		// Set an account to transfer to balance of the deleted account to
		SetTransferAccountID(hedera.AccountID{Account: 3}).
		SetTransactionMemo("go sdk example delete_account/main.go").
		Build(client)

	if err != nil {
		panic(err)
	}

	// Manually sign the transaction with the private key of the account to be deleted
	deleteTransaction = deleteTransaction.Sign(newKey)

	// Execute the transaction
	deleteTransactionID, err := deleteTransaction.Execute(client)

	if err != nil {
		panic(err)
	}

	deleteTransactionReceipt, err := deleteTransactionID.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("account delete transaction status: %v\n", deleteTransactionReceipt.Status)
}
