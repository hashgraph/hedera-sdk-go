package main

import (
	"fmt"
	"reflect"
	"unsafe"
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

	// Constructors exist for convenient files
	fileID := hedera.FileIDForAddressBook()
	// fileID := hedera.FileIDForFeeSchedule()
	// fileID := hedera.FileIDForExchangeRate()

	client.SetOperator(operatorAccountID, operatorPrivateKey)

	contents, err := hedera.NewFileContentsQuery().
		SetFileID(fileID).
		Execute(client)

	if err != nil {
		panic(err)
	}

	fmt.Printf("contents for file %v :\n", fileID)
	fmt.Print(string(contents))
	fmt.Println()
}
