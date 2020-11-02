package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
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

	if configOperatorID != "" && configOperatorKey != "" {
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

	// Constructors exist for convenient files
	fileID := hedera.FileIDForAddressBook()
	// fileID := hedera.FileIDForFeeSchedule()
	// fileID := hedera.FileIDForExchangeRate()

	contents, err := hedera.NewFileContentsQuery().
		SetFileID(fileID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("contents for file %v :\n", fileID)
	fmt.Print(string(contents))
	fmt.Println()
}
