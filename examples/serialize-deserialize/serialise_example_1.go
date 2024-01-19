package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	//client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	client, err = hedera.ClientForName("testnet")
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	//operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	operatorAccountID, err := hedera.AccountIDFromString("0.0.5698499")
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	//operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	operatorKey, err := hedera.PrivateKeyFromString("3030020100300706052b8104000a042204200e5f1866e11aa86cd9a00974d58a95b954f82469de7e8aa278ec6414ca752f8f")
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
