package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
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

	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(operatorAccountID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("balance = %v\n", balance)
}
