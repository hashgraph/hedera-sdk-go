package main

import (
	"encoding/json"
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

	info, err := hedera.NewAccountInfoQuery().
		SetAccountID(client.GetOperatorAccountID()).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing account info query")
		return
	}

	infoJSON, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		println(err.Error(), ": error marshaling to json")
		return
	}

	fmt.Printf("info for account %v :\n", client.GetOperatorAccountID())
	fmt.Print(string(infoJSON))
	fmt.Println()
}
