package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go"
)

func main() {
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	newKey, err := hedera.GenerateEd25519PrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Println("Manual signing example")
	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	client := hedera.ClientForTestnet()

	transaction, err := hedera.NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetInitialBalance(hedera.ZeroHbar).
		SetTransactionID(hedera.NewTransactionID(operatorAccountID)).
		SetTransactionMemo("sdk example create_account__with_manual_signing/main.go").
		Build(client)

	if err != nil {
		panic(err)
	}

	transactionID, err := transaction.
		Sign(operatorPrivateKey).
		Execute(client)

	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionID.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	newAccountId := transactionReceipt.GetAccountID()

	fmt.Printf("account = %v\n", newAccountId)
}
