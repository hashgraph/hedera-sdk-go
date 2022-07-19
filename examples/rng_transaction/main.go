package main

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	createResponse, err := hedera.NewPrngTransaction().
		// Set the range
		SetRange(12).
		Execute(client)
	if err != nil {
		println(err.Error(), ": error executing rng transaction")
		return
	}

	transactionRecord, err := createResponse.GetRecord(client)
	if err != nil {
		println(err.Error(), ": error getting receipt")
		return
	}

	if transactionRecord.PrngNumber != nil {
		println(err.Error(), ": error, pseudo-random number is nil")
		return
	}

	println("The pseudo-random number is:", *transactionRecord.PrngNumber)
}
