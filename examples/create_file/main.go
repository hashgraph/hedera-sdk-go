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

	transactionResponse, err := hedera.NewFileCreateTransaction().
		// A file is not implicitly owned by anyone, even the operator
		// But we do use operator's key for this one
		SetKeys(client.GetOperatorPublicKey()).
		// Initial contents of the file
		SetContents([]byte("Hello, World")).
		// Optional memo
		SetTransactionMemo("go sdk example create_file/main.go").
		// Set max transaction fee just in case we are required to pay more
		SetMaxTransactionFee(hedera.HbarFrom(8, hedera.HbarUnits.Hbar)).
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error creating file", err))
	}

	// Make sure the transaction went through
	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving file create transaction receipt", err))
	}

	// Get and then display the file ID from the receipt
	fmt.Printf("file = %v\n", *transactionReceipt.FileID)
}
