package main

import (
	"fmt"

	"github.com/arhtur007/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName("testnet")
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString("0.0.34195733")
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString("302e020100300506032b65700422042023acb1f99279eca4805b62027f34e560585d301986fb2530ef6cde1ac9177174")
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	fmt.Println("Crypto Transfer Example")

	fmt.Printf("Transferring 1 hbar from %v to 0.0.3\n", client.GetOperatorAccountID())

	transactionResponse, err := hedera.NewTransferTransaction().
		// Hbar has to be negated to denote we are taking out from that account
		AddHbarTransfer(client.GetOperatorAccountID(), hedera.NewHbar(-0.00000001)).
		// If the amount of these 2 transfers is not the same, the transaction will throw an error
		AddHbarTransfer(hedera.AccountID{Account: 34346808}, hedera.NewHbar(0.00000001)).
		SetTransactionMemo("go sdk example send_hbar/main.go").
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing transfer")
		return
	}

	// Retrieve the receipt to make sure the transaction went through
	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		println(err.Error(), ": error retrieving transfer receipt")
		return
	}

	fmt.Printf("crypto transfer status: %v\n", transactionReceipt.Status)
}
