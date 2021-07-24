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

	fmt.Println("Creating a file to delete:")

	// first create a file

	transactionResponse, err := hedera.NewFileCreateTransaction().
		SetContents([]byte("The quick brown fox jumps over the lazy dog")).
		SetKeys(client.GetOperatorPublicKey()).
		SetTransactionMemo("go sdk example delete_file/main.go").
		SetMaxTransactionFee(hedera.HbarFrom(8, hedera.HbarUnits.Hbar)).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating file")
		return
	}

	receipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving file creation receipt")
		return
	}

	newFileID := *receipt.FileID

	fmt.Printf("file = %v\n", newFileID)
	fmt.Println("deleting created file")

	// To delete a file you must do the following:
	deleteTransactionID, err := hedera.NewFileDeleteTransaction().
		SetFileID(newFileID).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error deleting file")
		return
	}

	deleteTransactionReceipt, err := deleteTransactionID.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving file deletion receipt")
		return
	}

	fmt.Printf("file delete transaction status: %v\n", deleteTransactionReceipt.Status)

	// querying for file info on a deleted file will result in FILE_DELETED
	fileInfo, err := hedera.NewFileInfoQuery().
		SetFileID(newFileID).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing file info query")
		return
	}

	fmt.Printf("file %v was deleted: %v\n", newFileID, fileInfo.IsDeleted)
}
