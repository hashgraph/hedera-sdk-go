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

	client.SetOperator(
		// Operator Account ID
		hedera.AccountID{Account: 2},
		// Operator Private Key
		operatorPrivateKey,
	)

	transactionID, err := hedera.NewFileCreateTransaction().
		AddKey(operatorPrivateKey.PublicKey()).
		SetContents([]byte{1, 2, 3, 4}).
		SetTransactionMemo("sdk example create_file/main.go").
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.Receipt(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("file = %v\n", transactionReceipt.FileID())
}
