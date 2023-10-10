package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(fmt.Sprintf("%v : error creating client", err))
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to AccountID", err))
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(fmt.Sprintf("%v : error converting string to PrivateKey", err))
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	fileQuery := hedera.NewFileContentsQuery().
		// Set the file ID for address book which is 0.0.102
		SetFileID(hedera.FileIDForAddressBook())

	println("the network that address book is for:", client.GetNetworkName().String())

	cost, err := fileQuery.GetCost(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error getting file contents query cost", err))
	}

	println("file contents cost:", cost.String())

	// Have to always set the cost a little bigger, otherwise it is possible it won't go through
	fileQuery.SetMaxQueryPayment(hedera.NewHbar(1))

	// Execute the file content query
	contents, err := fileQuery.Execute(client)
	if err != nil {
		panic(fmt.Sprintf("%v : error executing file contents query", err))
	}

	fileByte, err := os.OpenFile("address-book-byte.pb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("%v : error opening address-book-byte.pb", err))
	}

	fileString, err := os.OpenFile("address-book-string.pb", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("%v : error opening address-book-string.pb", err))
	}

	// Write the contents (string([]byte)) into the string file
	leng, err := fileString.WriteString(string(contents))
	if err != nil {
		panic(fmt.Sprintf("%v : error writing contents to file", err))
	}
	// Write the contents ([]byte) into the byte file
	_, err = fileByte.Write(contents)
	if err != nil {
		panic(fmt.Sprintf("%v : error writing contents to file", err))
	}

	temp := make([]byte, leng)

	_, err = fileString.Read(temp)

	// Close the files
	err = fileString.Close()
	if err != nil {
		panic(fmt.Sprintf("%v : error closing the file", err))
	}
	err = fileByte.Close()
	if err != nil {
		panic(fmt.Sprintf("%v : error closing the file", err))
	}
}
