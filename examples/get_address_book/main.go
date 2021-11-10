package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	client.SetOperator(operatorAccountID, operatorKey)

	fileQuery := hedera.NewFileContentsQuery().
		SetFileID(hedera.FileIDForAddressBook())

	println("file contents cost: ", client.GetOperatorAccountID().String())

	cost, err := fileQuery.GetCost(client)
	if err != nil {
		println(err.Error(), ": error getting file contents query cost")
		return
	}

	println("file contents cost: ", cost.String())

	fileQuery.SetMaxQueryPayment(hedera.NewHbar(1))

	contents, err := fileQuery.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing file contents query")
		return
	}

	file, err := os.OpenFile("address-book.services.bin", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println(err.Error(), ": error opening address-book.services.bin")
		return
	}

	err = file.Truncate(0)
	if err != nil {
		println(err.Error(), ": error truncating file")
		return
	}

	_, err = fmt.Fprintf(file, "%d", contents)

}
