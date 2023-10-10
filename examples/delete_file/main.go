package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	// Generate the key to be used with the new file
	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(fmt.Sprintf("%v : error generating PrivateKey", err))
	}

	fmt.Println("Creating a file to delete:")

	// First create a file
	freezeTransaction, err := hedera.NewFileCreateTransaction().
		// Mock contents
		SetContents([]byte("The quick brown fox jumps over the lazy dog")).
		// All keys at the top level of a key list must sign to create or modify the file. Any one of
		// the keys at the top level key list can sign to delete the file.
		SetKeys(newKey.PublicKey()).
		SetTransactionMemo("go sdk example delete_file/main.go").
		SetMaxTransactionFee(hedera.HbarFrom(8, hedera.HbarUnits.Hbar)).
		FreezeWith(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error freezing transaction", err))
	}
	transactionResponse, err := freezeTransaction.Sign(newKey).Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error creating file", err))
	}

	// Get the receipt to make sure transaction went through
	receipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving file creation receipt", err))
	}

	// Retrieve file ID from the receipt
	newFileID := *receipt.FileID

	fmt.Printf("file = %v\n", newFileID)
	fmt.Println("deleting created file")

	// To delete a file you must do the following:
	deleteTransaction, err := hedera.NewFileDeleteTransaction().
		// Set file ID
		SetFileID(newFileID).
		FreezeWith(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error freezing file delete transaction", err))
	}

	// Sign with the key we used to create the file
	deleteTransaction.Sign(newKey)

	// Execute the file delete transaction
	deleteTransactionResponse, err := deleteTransaction.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing file delete transaction", err))
	}

	// Check that it went through
	deleteTransactionReceipt, err := deleteTransactionResponse.GetReceipt(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving file deletion receipt", err))
	}

	fmt.Printf("file delete transaction status: %v\n", deleteTransactionReceipt.Status)

	// Querying for file info on a deleted file will result in FILE_DELETED
	// Good way to check if file was actually deleted
	fileInfo, err := hedera.NewFileInfoQuery().
		// Only file ID required
		SetFileID(newFileID).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error executing file info query", err))
	}

	fmt.Printf("file %v was deleted: %v\n", newFileID, fileInfo.IsDeleted)
}
