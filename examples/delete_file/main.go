package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
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

	fmt.Println("Creating a file to delete:")

	// first create a file

	transactionResponse, err := hedera.NewFileCreateTransaction().
		SetContents([]byte("The quick brown fox jumps over the lazy dog")).
		SetKeys(client.GetOperatorPublicKey()).
		SetTransactionMemo("go sdk example delete_file/main.go").
		SetMaxTransactionFee(hedera.HbarFrom(8, hedera.HbarUnits.Hbar)).
		Execute(client)

	if err != nil {
		panic(err)
	}

	receipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	newFileID := *receipt.FileID

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

	fmt.Printf("file %v was deleted: %v\n", newFileID, fileInfo.IsDeleted)
}
