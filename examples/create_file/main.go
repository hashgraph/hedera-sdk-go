package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	var client *hedera.Client
	var err error

	if os.Getenv("HEDERA_NETWORK") == "previewnet" {
		client = hedera.ClientForPreviewnet()
	} else {
		client, err = hedera.ClientFromConfigFile(os.Getenv("CONFIG_FILE"))

		if err != nil {
			println("not error", err.Error())
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")
	var operatorKey hedera.PrivateKey

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err = hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	transactionResponse, err := hedera.NewFileCreateTransaction().
		// A file is not implicitly owned by anyone, even the operator
		SetKeys(client.GetOperatorKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk example create_file/main.go").
		SetMaxTransactionFee(hedera.HbarFrom(8, hedera.HbarUnits.Hbar)).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("file = %v\n", *transactionReceipt.FileID)
}
