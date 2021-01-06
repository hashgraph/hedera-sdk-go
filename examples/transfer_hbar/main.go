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

	fmt.Println("Crypto Transfer Example")

	fmt.Printf("Transferring 1 hbar from %v to 0.0.3\n", client.GetOperatorAccountID())

	transactionResponse, err := hedera.NewTransferTransaction().
		AddHbarTransfer(client.GetOperatorAccountID(), hedera.NewHbar(-1)).
		AddHbarTransfer(hedera.AccountID{Account: 3}, hedera.NewHbar(1)).
		SetTransactionMemo("go sdk example send_hbar/main.go").
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing transfer")
		return
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		println(err.Error(), ": error retrieving transfer receipt")
		return
	}

	fmt.Printf("crypto transfer status: %v\n", transactionReceipt.Status)
}
