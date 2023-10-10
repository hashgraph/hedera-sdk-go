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

	createResponse, err := hedera.NewPrngTransaction().
		// Set the range
		SetRange(12).
		Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing rng transaction", err))
	}

	transactionRecord, err := createResponse.GetRecord(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting receipt", err))
	}
	fmt.Printf("transactionRecord: %v\n", transactionRecord)
	if transactionRecord.PrngNumber == nil {
		panic(fmt.Sprintf("%v : error, pseudo-random number is nil", err))
	}

	println("The pseudo-random number is:", *transactionRecord.PrngNumber)
}
