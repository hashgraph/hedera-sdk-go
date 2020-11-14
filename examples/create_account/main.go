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
			client = hedera.ClientForTestnet()
		}
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" && client.GetOperatorPublicKey().Bytes() == nil {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	transactionResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetTransactionMemo("go sdk example create_account/main.go").
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)
}
