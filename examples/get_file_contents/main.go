package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
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
			println(err.Error(), ": error setting up client from config file")
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" && client.GetOperatorPublicKey().Bytes() == nil {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			println(err.Error(), ": error converting string to AccountID")
			return
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			println(err.Error(), ": error converting string to PrivateKey")
			return
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
		println(err.Error(), ": error executing file contents query")
		return
	}

	fmt.Printf("contents for file %v :\n", fileID.String())
	fmt.Print(string(contents))
	fmt.Println()
}
