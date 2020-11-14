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

	if configOperatorID != "" && configOperatorKey != "" && client.GetOperatorPublicKey().Bytes() == nil {
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

	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Println("Manual signing example")
	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	transaction, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(hedera.ZeroHbar).
		SetTransactionID(hedera.TransactionIDGenerate(client.GetOperatorAccountID())).
		SetTransactionMemo("sdk example create_account__with_manual_signing/main.go").
		FreezeWith(client)

	if err != nil {
		panic(err)
	}

	transactionResponse, err := transaction.
		Sign(operatorKey).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	newAccountId := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountId)
}
