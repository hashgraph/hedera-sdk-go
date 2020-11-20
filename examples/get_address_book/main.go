package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
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


	fileQuery := hedera.NewFileContentsQuery().
		SetFileID(hedera.FileIDForAddressBook())

	println("file contents cost: ", client.GetOperatorAccountID().String())

	cost, err := fileQuery.GetCost(client)
	if err != nil {
		panic(err)
	}

	println("file contents cost: ", cost.String())

	fileQuery.SetMaxQueryPayment(hedera.NewHbar(1))

	contents, err := fileQuery.Execute(client)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("address-book.proto.bin", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	err = file.Truncate(0)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(file, "%d", contents)

}
