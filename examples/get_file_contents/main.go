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

	// Constructors exist for convenient files
	//fileID := hedera.FileIDForAddressBook()
	// fileID := hedera.FileIDForFeeSchedule()
	fileID := hedera.FileIDForExchangeRate()

	contents, err := hedera.NewFileContentsQuery().
		// Set the file ID for whatever file you need
		SetFileID(fileID).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error executing file contents query", err))
	}

	exchangeRate, err := hedera.ExchangeRateFromBytes(contents)
	if err != nil {
		panic(fmt.Sprintf("%v : error converting contents to exchange rate", err))
	}

	fmt.Printf("Contents for file %v :\n", fileID.String())
	fmt.Print(exchangeRate.String())
	fmt.Println()
}
