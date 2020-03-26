package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
)

func main() {
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Creating a file to delete:")
	client := hedera.ClientForTestnet().
		SetOperator(operatorAccountID, operatorPrivateKey)

	// first create a file

	transactionID, err := hedera.NewFileCreateTransaction().
		SetContents([]byte("The quick brown fox jumps over the lazy dog")).
		SetTransactionMemo("go sdk example delete_file/main.go").
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	newFileID := transactionReceipt.GetFileID()

	fmt.Printf("file = %v\n", newFileID)
	fmt.Println("deleting created file")

	// To delete a file you must do the following:
	deleteTransactionID, err := hedera.NewFileDeleteTransaction().
		SetFileID(newFileID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	deleteTransactionReceipt, err := deleteTransactionID.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("file delete transaction status: %v\n", deleteTransactionReceipt.Status)

	// querying for file info on a deleted file will result in FILE_DELETED
	fileInfo, err := hedera.NewFileInfoQuery().
		SetFileID(newFileID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("file %v was deleted: %\nv", newFileID, fileInfo.IsDeleted)
}
