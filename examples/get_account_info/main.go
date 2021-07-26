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

	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	client.SetOperator(operatorAccountID, operatorKey)

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
