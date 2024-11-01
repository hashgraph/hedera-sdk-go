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

	// Generate new key to use with new account
	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	resp, err := hedera.NewAccountCreateTransaction().SetKey(newKey).Execute(client)
	receipt, err := resp.GetReceipt(client)
	newAccountId := *receipt.AccountID

	bytes, err := hedera.NewTransferTransaction().AddHbarTransfer(operatorAccountID, hedera.NewHbar(1)).
		ToBytes()

	if err != nil {
		panic(err)
	}

	txFromBytes, err := hedera.TransactionFromBytes(bytes)

	transaction := txFromBytes.(hedera.TransferTransaction)
	_, err = transaction.AddHbarTransfer(newAccountId, hedera.NewHbar(-1)).SignWithOperator(client)

	_, err = transaction.Execute(client)
	if err != nil {
		panic(err)
	}
	// Get the `AccountInfo` on the new account and show it is a hollow account by not having a public key
	info, err := hedera.NewAccountInfoQuery().SetAccountID(newAccountId).Execute(client)
	fmt.Println("Balance of new account: ", info.Balance)
}
