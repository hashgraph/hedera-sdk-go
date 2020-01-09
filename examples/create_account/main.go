package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
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

	client.SetOperator(
		// Operator Account ID
		operatorAccountID,
		// Operator Private Key
		operatorPrivateKey,
	)

	newKey, err := hedera.GenerateEd25519PrivateKey()

	if err != nil {
		panic(err)
	}

	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	transactionID, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(1000).
		SetTransactionMemo("sdk example create_account/main.go").
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.Receipt(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("account = %v\n", transactionReceipt.AccountID())
}
