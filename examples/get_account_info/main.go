package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
)

func main() {
	client := hedera.ClientForTestnet()

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	client.SetOperator(operatorAccountID, operatorPrivateKey)

	info, err := hedera.NewAccountInfoQuery().
		SetAccountID(operatorAccountID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	infoJSON, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("info for account %v :\n", operatorAccountID)
	fmt.Print(string(infoJSON))
	fmt.Println()
}
