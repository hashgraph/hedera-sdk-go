package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
)

func main() {
	client := hedera.NewClient(map[string]hedera.AccountID{
		"0.testnet.hedera.com:50211": {Account: 3},
	})

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
		SetAccountID(oepratorAccountID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("balance = %v tħ\n", balance)
}
