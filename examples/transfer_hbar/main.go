package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
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

	fmt.Println("Crypto Transfer Example")

	client := hedera.ClientForTestnet().
		SetOperator(operatorAccountID, operatorPrivateKey)

	fmt.Printf("Transferring 1 hbar from %v to 0.0.3\n", operatorAccountID)

	transactionID, err := hedera.NewCryptoTransferTransaction().
		AddSender(operatorAccountID, hedera.NewHbar(1)).
		AddRecipient(hedera.AccountID{Account: 3}, hedera.NewHbar(1)).
		SetTransactionMemo("go sdk example send_hbar/main.go").
		Execute(client)

	if err != nil {
		panic(err)
	}

	receipt, err := transactionID.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("crypto transfer status: %v\n", receipt.Status)
}
