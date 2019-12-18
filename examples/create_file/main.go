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

	client.SetOperator(
		// Operator Account ID
		hedera.AccountID{Account: 2},
		// Operator Private Key
		operatorPrivateKey,
	)

	transactionId, err := hedera.NewFileCreateTransaction().
		AddKey(operatorPrivateKey.PublicKey()).
		SetContents([]byte{1, 2, 3, 4}).
		SetMemo("sdk example create_file/main.go").
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionId.Receipt(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("file = %v\n", transactionReceipt.FileID)
}
