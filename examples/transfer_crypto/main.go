package main

import (
	"fmt"
	"os"

	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

func main() {
	var client *hiero.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hiero.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hiero.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hiero.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	fmt.Println("Crypto Transfer Example")

	fmt.Printf("Transferring 1 hbar from %v to 0.0.3\n", client.GetOperatorAccountID())

	transactionResponse, err := hiero.NewTransferTransaction().
		// Hbar has to be negated to denote we are taking out from that account
		AddHbarTransfer(client.GetOperatorAccountID(), hiero.NewHbar(-1)).
		// If the amount of these 2 transfers is not the same, the transaction will throw an error
		AddHbarTransfer(hiero.AccountID{Account: 3}, hiero.NewHbar(1)).
		SetTransactionMemo("go sdk example send_hbar/main.go").
		Execute(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error executing transfer", err))
	}

	// Retrieve the receipt to make sure the transaction went through
	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		panic(fmt.Sprintf("%v : error retrieving transfer receipt", err))
	}

	fmt.Printf("crypto transfer status: %v\n", transactionReceipt.Status)
}
