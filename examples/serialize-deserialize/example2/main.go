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

	// Generate new key to use with new account
	newKey, err := hiero.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	resp, err := hiero.NewAccountCreateTransaction().SetKey(newKey).Execute(client)
	receipt, err := resp.GetReceipt(client)
	newAccountId := *receipt.AccountID

	// Prepare and sign the tx and send it to be signed by another actor
	fmt.Println("Creating a transfer transaction, signing it with operator and serializing it to bytes...")
	bytes, err := hiero.NewTransferTransaction().AddHbarTransfer(operatorAccountID, hiero.NewHbar(1)).AddHbarTransfer(newAccountId, hiero.NewHbar(-1)).
		Sign(operatorKey).ToBytes()

	FromBytes, err := hiero.TransactionFromBytes(bytes)
	if err != nil {
		panic(err)
	}
	txFromBytes := FromBytes.(hiero.TransferTransaction)
	// New Account add his sign and execute the tx:
	fmt.Println("Signing deserialized transaction with `newAccount` private key and executing it...")
	executed, err := txFromBytes.Sign(newKey).SetMaxTransactionFee(hiero.NewHbar(2)).Execute(client)
	if err != nil {
		panic(err)
	}
	receipt, err = executed.GetReceipt(client)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tx successfully executed. Here is receipt:", receipt)
}
