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

	client.SetOperator(hedera.AccountID{Account: 2}, operatorPrivateKey)

	balance, err := hedera.NewAccountBalanceQuery().
		SetAccountID(hedera.AccountID{Account: 2}).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("balance = %v tħ\n", balance)
}
