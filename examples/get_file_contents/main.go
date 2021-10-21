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

	// Constructors exist for convenient files
	fileID := hedera.FileIDForAddressBook()
	// fileID := hedera.FileIDForFeeSchedule()
	// fileID := hedera.FileIDForExchangeRate()

	contents, err := hedera.NewFileContentsQuery().
		SetFileID(fileID).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing file contents query")
		return
	}

	fmt.Printf("contents for file %v :\n", fileID.String())
	fmt.Print(string(contents))
	fmt.Println()
}
