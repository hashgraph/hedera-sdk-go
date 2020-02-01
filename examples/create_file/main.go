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

	transactionID, err := hedera.NewFileCreateTransaction().
		// A file is not implicitly owned by anyone, even the operator
		AddKey(operatorPrivateKey.PublicKey()).
		SetContents([]byte("Hello, World")).
		SetTransactionMemo("go sdk example create_file/main.go").
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("file = %v\n", transactionReceipt.GetFileID())
}
